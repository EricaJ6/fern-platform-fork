// Package dataloader provides efficient batch loading for GraphQL resolvers
package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/guidewire-oss/fern-platform/pkg/database"
	"gorm.io/gorm"
)

// contextKey is used for storing dataloaders in context
type contextKey string

const (
	projectLoaderKey  contextKey = "projectLoader"
	testRunLoaderKey  contextKey = "testRunLoader"
	suiteRunLoaderKey contextKey = "suiteRunLoader"
	specRunLoaderKey  contextKey = "specRunLoader"
	tagLoaderKey      contextKey = "tagLoader"
	userLoaderKey     contextKey = "userLoader"
)

// Loaders holds all the dataloaders
type Loaders struct {
	ProjectByID  *dataloader.Loader[string, *database.ProjectDetails]
	TestRunByID  *dataloader.Loader[string, *database.TestRun]
	SuiteRunByID *dataloader.Loader[string, *database.SuiteRun]
	SpecRunByID  *dataloader.Loader[string, *database.SpecRun]
	TagByID      *dataloader.Loader[string, *database.Tag]
	UserByID     *dataloader.Loader[string, *database.User]

	// Batch loaders for relationships
	SuiteRunsByTestRunID *dataloader.Loader[string, []*database.SuiteRun]
	SpecRunsBySuiteRunID *dataloader.Loader[string, []*database.SpecRun]
	TagsByTestRunID      *dataloader.Loader[string, []*database.Tag]
}

// NewLoaders creates a new set of dataloaders
func NewLoaders(db *gorm.DB) *Loaders {
	return &Loaders{
		ProjectByID:  createProjectLoader(db),
		TestRunByID:  createTestRunLoader(db),
		SuiteRunByID: createSuiteRunLoader(db),
		SpecRunByID:  createSpecRunLoader(db),
		TagByID:      createTagLoader(db),
		UserByID:     createUserLoader(db),

		SuiteRunsByTestRunID: createSuiteRunsByTestRunLoader(db),
		SpecRunsBySuiteRunID: createSpecRunsBySuiteRunLoader(db),
		TagsByTestRunID:      createTagsByTestRunLoader(db),
	}
}

// Middleware adds dataloaders to the context
func Middleware(loaders *Loaders) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, projectLoaderKey, loaders.ProjectByID)
			ctx = context.WithValue(ctx, testRunLoaderKey, loaders.TestRunByID)
			ctx = context.WithValue(ctx, suiteRunLoaderKey, loaders.SuiteRunByID)
			ctx = context.WithValue(ctx, specRunLoaderKey, loaders.SpecRunByID)
			ctx = context.WithValue(ctx, tagLoaderKey, loaders.TagByID)
			ctx = context.WithValue(ctx, userLoaderKey, loaders.UserByID)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetProjectLoader returns the project loader from context
func GetProjectLoader(ctx context.Context) *dataloader.Loader[string, *database.ProjectDetails] {
	return ctx.Value(projectLoaderKey).(*dataloader.Loader[string, *database.ProjectDetails])
}

// GetTestRunLoader returns the test run loader from context
func GetTestRunLoader(ctx context.Context) *dataloader.Loader[string, *database.TestRun] {
	return ctx.Value(testRunLoaderKey).(*dataloader.Loader[string, *database.TestRun])
}

// Batch loading functions

func createProjectLoader(db *gorm.DB) *dataloader.Loader[string, *database.ProjectDetails] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.ProjectDetails] {
		// Create a map for quick lookup
		projectMap := make(map[string]*database.ProjectDetails)

		// Batch query
		var projects []database.ProjectDetails
		if err := db.Where("id IN ?", keys).Find(&projects).Error; err != nil {
			// Return error for all keys
			results := make([]*dataloader.Result[*database.ProjectDetails], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.ProjectDetails]{Error: err}
			}
			return results
		}

		// Build map
		for i := range projects {
			projectMap[fmt.Sprintf("%d", projects[i].ID)] = &projects[i]
		}

		// Build results in the same order as keys
		results := make([]*dataloader.Result[*database.ProjectDetails], len(keys))
		for i, key := range keys {
			if project, ok := projectMap[key]; ok {
				results[i] = &dataloader.Result[*database.ProjectDetails]{Data: project}
			} else {
				results[i] = &dataloader.Result[*database.ProjectDetails]{Error: fmt.Errorf("project not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.ProjectDetails]{}))
}

func createTestRunLoader(db *gorm.DB) *dataloader.Loader[string, *database.TestRun] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.TestRun] {
		testRunMap := make(map[string]*database.TestRun)

		var testRuns []database.TestRun
		if err := db.Where("id IN ?", keys).Find(&testRuns).Error; err != nil {
			results := make([]*dataloader.Result[*database.TestRun], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.TestRun]{Error: err}
			}
			return results
		}

		for i := range testRuns {
			testRunMap[fmt.Sprintf("%d", testRuns[i].ID)] = &testRuns[i]
		}

		results := make([]*dataloader.Result[*database.TestRun], len(keys))
		for i, key := range keys {
			if testRun, ok := testRunMap[key]; ok {
				results[i] = &dataloader.Result[*database.TestRun]{Data: testRun}
			} else {
				results[i] = &dataloader.Result[*database.TestRun]{Error: fmt.Errorf("test run not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.TestRun]{}))
}

func createSuiteRunLoader(db *gorm.DB) *dataloader.Loader[string, *database.SuiteRun] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.SuiteRun] {
		suiteRunMap := make(map[string]*database.SuiteRun)

		var suiteRuns []database.SuiteRun
		if err := db.Where("id IN ?", keys).Find(&suiteRuns).Error; err != nil {
			results := make([]*dataloader.Result[*database.SuiteRun], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.SuiteRun]{Error: err}
			}
			return results
		}

		for i := range suiteRuns {
			suiteRunMap[fmt.Sprintf("%d", suiteRuns[i].ID)] = &suiteRuns[i]
		}

		results := make([]*dataloader.Result[*database.SuiteRun], len(keys))
		for i, key := range keys {
			if suiteRun, ok := suiteRunMap[key]; ok {
				results[i] = &dataloader.Result[*database.SuiteRun]{Data: suiteRun}
			} else {
				results[i] = &dataloader.Result[*database.SuiteRun]{Error: fmt.Errorf("suite run not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.SuiteRun]{}))
}

func createSpecRunLoader(db *gorm.DB) *dataloader.Loader[string, *database.SpecRun] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.SpecRun] {
		specRunMap := make(map[string]*database.SpecRun)

		var specRuns []database.SpecRun
		if err := db.Where("id IN ?", keys).Find(&specRuns).Error; err != nil {
			results := make([]*dataloader.Result[*database.SpecRun], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.SpecRun]{Error: err}
			}
			return results
		}

		for i := range specRuns {
			specRunMap[fmt.Sprintf("%d", specRuns[i].ID)] = &specRuns[i]
		}

		results := make([]*dataloader.Result[*database.SpecRun], len(keys))
		for i, key := range keys {
			if specRun, ok := specRunMap[key]; ok {
				results[i] = &dataloader.Result[*database.SpecRun]{Data: specRun}
			} else {
				results[i] = &dataloader.Result[*database.SpecRun]{Error: fmt.Errorf("spec run not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.SpecRun]{}))
}

func createTagLoader(db *gorm.DB) *dataloader.Loader[string, *database.Tag] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.Tag] {
		tagMap := make(map[string]*database.Tag)

		var tags []database.Tag
		if err := db.Where("id IN ?", keys).Find(&tags).Error; err != nil {
			results := make([]*dataloader.Result[*database.Tag], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.Tag]{Error: err}
			}
			return results
		}

		for i := range tags {
			tagMap[fmt.Sprintf("%d", tags[i].ID)] = &tags[i]
		}

		results := make([]*dataloader.Result[*database.Tag], len(keys))
		for i, key := range keys {
			if tag, ok := tagMap[key]; ok {
				results[i] = &dataloader.Result[*database.Tag]{Data: tag}
			} else {
				results[i] = &dataloader.Result[*database.Tag]{Error: fmt.Errorf("tag not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.Tag]{}))
}

func createUserLoader(db *gorm.DB) *dataloader.Loader[string, *database.User] {
	batchFn := func(ctx context.Context, keys []string) []*dataloader.Result[*database.User] {
		userMap := make(map[string]*database.User)

		var users []database.User
		if err := db.Where("user_id IN ?", keys).Find(&users).Error; err != nil {
			results := make([]*dataloader.Result[*database.User], len(keys))
			for i := range results {
				results[i] = &dataloader.Result[*database.User]{Error: err}
			}
			return results
		}

		for i := range users {
			userMap[users[i].UserID] = &users[i]
		}

		results := make([]*dataloader.Result[*database.User], len(keys))
		for i, key := range keys {
			if user, ok := userMap[key]; ok {
				results[i] = &dataloader.Result[*database.User]{Data: user}
			} else {
				results[i] = &dataloader.Result[*database.User]{Error: fmt.Errorf("user not found: %s", key)}
			}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(&dataloader.InMemoryCache[string, *database.User]{}))
}

// Relationship loaders

func createSuiteRunsByTestRunLoader(db *gorm.DB) *dataloader.Loader[string, []*database.SuiteRun] {
	batchFn := func(ctx context.Context, testRunIDs []string) []*dataloader.Result[[]*database.SuiteRun] {
		// Convert string IDs to integers
		intIDs := make([]int, len(testRunIDs))
		for i, id := range testRunIDs {
			parsedID, err := strconv.Atoi(id)
			if err != nil {
				results := make([]*dataloader.Result[[]*database.SuiteRun], len(testRunIDs))
				for j := range results {
					results[j] = &dataloader.Result[[]*database.SuiteRun]{Error: fmt.Errorf("invalid test run ID: %s", testRunIDs[j])}
				}
				return results
			}
			intIDs[i] = parsedID
		}

		// Query all suite runs for all test run IDs at once
		var suiteRuns []database.SuiteRun
		if err := db.Where("test_run_id IN ?", intIDs).
			Order("start_time ASC").
			Find(&suiteRuns).Error; err != nil {
			results := make([]*dataloader.Result[[]*database.SuiteRun], len(testRunIDs))
			for i := range results {
				results[i] = &dataloader.Result[[]*database.SuiteRun]{Error: err}
			}
			return results
		}

		// Group by test run ID
		suiteRunsByTestRun := make(map[string][]*database.SuiteRun)
		for i := range suiteRuns {
			testRunID := fmt.Sprintf("%d", suiteRuns[i].TestRunID)
			suiteRunsByTestRun[testRunID] = append(suiteRunsByTestRun[testRunID], &suiteRuns[i])
		}

		// Build results
		results := make([]*dataloader.Result[[]*database.SuiteRun], len(testRunIDs))
		for i, testRunID := range testRunIDs {
			suites := suiteRunsByTestRun[testRunID]
			if suites == nil {
				suites = []*database.SuiteRun{} // Return empty slice instead of nil
			}
			results[i] = &dataloader.Result[[]*database.SuiteRun]{Data: suites}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn,
		dataloader.WithCache(&dataloader.InMemoryCache[string, []*database.SuiteRun]{}),
		dataloader.WithBatchCapacity[string, []*database.SuiteRun](100),
		dataloader.WithWait[string, []*database.SuiteRun](2*time.Millisecond))
}

func createSpecRunsBySuiteRunLoader(db *gorm.DB) *dataloader.Loader[string, []*database.SpecRun] {
	batchFn := func(ctx context.Context, suiteRunIDs []string) []*dataloader.Result[[]*database.SpecRun] {
		// Convert string IDs to integers
		intIDs := make([]int, len(suiteRunIDs))
		for i, id := range suiteRunIDs {
			parsedID, err := strconv.Atoi(id)
			if err != nil {
				results := make([]*dataloader.Result[[]*database.SpecRun], len(suiteRunIDs))
				for j := range results {
					results[j] = &dataloader.Result[[]*database.SpecRun]{Error: fmt.Errorf("invalid suite run ID: %s", suiteRunIDs[j])}
				}
				return results
			}
			intIDs[i] = parsedID
		}

		var specRuns []database.SpecRun
		if err := db.Where("suite_run_id IN ?", intIDs).
			Order("start_time ASC").
			Find(&specRuns).Error; err != nil {
			results := make([]*dataloader.Result[[]*database.SpecRun], len(suiteRunIDs))
			for i := range results {
				results[i] = &dataloader.Result[[]*database.SpecRun]{Error: err}
			}
			return results
		}

		specRunsBySuite := make(map[string][]*database.SpecRun)
		for i := range specRuns {
			suiteRunID := fmt.Sprintf("%d", specRuns[i].SuiteRunID)
			specRunsBySuite[suiteRunID] = append(specRunsBySuite[suiteRunID], &specRuns[i])
		}

		results := make([]*dataloader.Result[[]*database.SpecRun], len(suiteRunIDs))
		for i, suiteRunID := range suiteRunIDs {
			specs := specRunsBySuite[suiteRunID]
			if specs == nil {
				specs = []*database.SpecRun{}
			}
			results[i] = &dataloader.Result[[]*database.SpecRun]{Data: specs}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn,
		dataloader.WithCache(&dataloader.InMemoryCache[string, []*database.SpecRun]{}),
		dataloader.WithBatchCapacity[string, []*database.SpecRun](100),
		dataloader.WithWait[string, []*database.SpecRun](2*time.Millisecond))
}

func createTagsByTestRunLoader(db *gorm.DB) *dataloader.Loader[string, []*database.Tag] {
	batchFn := func(ctx context.Context, testRunIDs []string) []*dataloader.Result[[]*database.Tag] {
		// Query tag associations
		type tagAssoc struct {
			TestRunID uint
			TagID     uint
		}

		var associations []tagAssoc
		if err := db.Table("test_run_tags").
			Select("test_run_id, tag_id").
			Where("test_run_id IN ?", testRunIDs).
			Scan(&associations).Error; err != nil {
			results := make([]*dataloader.Result[[]*database.Tag], len(testRunIDs))
			for i := range results {
				results[i] = &dataloader.Result[[]*database.Tag]{Error: err}
			}
			return results
		}

		// Get unique tag IDs
		tagIDMap := make(map[uint]bool)
		for _, assoc := range associations {
			tagIDMap[assoc.TagID] = true
		}

		tagIDs := make([]uint, 0, len(tagIDMap))
		for id := range tagIDMap {
			tagIDs = append(tagIDs, id)
		}

		// Fetch all tags
		var tags []database.Tag
		if len(tagIDs) > 0 {
			if err := db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
				results := make([]*dataloader.Result[[]*database.Tag], len(testRunIDs))
				for i := range results {
					results[i] = &dataloader.Result[[]*database.Tag]{Error: err}
				}
				return results
			}
		}

		// Build tag map
		tagMap := make(map[uint]*database.Tag)
		for i := range tags {
			tagMap[tags[i].ID] = &tags[i]
		}

		// Group tags by test run
		tagsByTestRun := make(map[string][]*database.Tag)
		for _, assoc := range associations {
			if tag, ok := tagMap[assoc.TagID]; ok {
				tagsByTestRun[fmt.Sprintf("%d", assoc.TestRunID)] = append(tagsByTestRun[fmt.Sprintf("%d", assoc.TestRunID)], tag)
			}
		}

		// Build results
		results := make([]*dataloader.Result[[]*database.Tag], len(testRunIDs))
		for i, testRunID := range testRunIDs {
			tags := tagsByTestRun[testRunID]
			if tags == nil {
				tags = []*database.Tag{}
			}
			results[i] = &dataloader.Result[[]*database.Tag]{Data: tags}
		}

		return results
	}

	return dataloader.NewBatchedLoader(batchFn,
		dataloader.WithCache(&dataloader.InMemoryCache[string, []*database.Tag]{}),
		dataloader.WithBatchCapacity[string, []*database.Tag](100),
		dataloader.WithWait[string, []*database.Tag](2*time.Millisecond),
	)
}
