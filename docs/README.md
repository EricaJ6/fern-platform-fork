# Fern Platform Documentation

<div align="center">
  <img src="https://github.com/guidewire-oss/fern-reporter/blob/main/docs/images/logo-color.png" alt="Fern Platform" width="200"/>
  
  **A unified platform for test reporting, analysis, and AI-powered insights**
</div>

## 📚 Documentation Navigation

Choose your path based on your role and goals:

### 🎯 **For Product Managers & Stakeholders**
- [**Product Overview**](product/overview.md) - What Fern Platform does and business value
- [**UI Enhancements**](UI_ENHANCEMENTS.md) - Modern dashboard and visualization features
- [**Architecture Overview**](ARCHITECTURE.md) - High-level system design

### 🛠️ **For Developers & Operators**
- [**Developer Guide**](developers/README.md) - Complete developer documentation hub
- [**Quick Start Guide**](developers/quick-start.md) - Get running in 15 minutes
- [**Local k3d Setup**](installation/local-k3d.md) - Complete local Kubernetes setup
- [**REST API Reference**](developers/api-reference.md) - REST endpoints documentation
- [**GraphQL API**](graphql-api.md) - GraphQL schema and queries
- [**OAuth Configuration**](configuration/oauth.md) - Authentication setup

### 🏗️ **For Platform Engineers & Architects**
- [**Architecture Document**](ARCHITECTURE.md) - System design and technical decisions
- [**RFCs**](rfc/) - Technical proposals and future plans
  - [Platform Consolidation](rfc/rfc-001-platform-consolidation-and-architecture-evolution.md)
  - [LLM Integration](rfc/rfc-002-llm-provider-integration-and-ai-intelligence-architecture.md)
  - [Requirements Traceability](rfc/rfc-003-requirements-traceability-and-test-coverage-intelligence.md)

### 🚀 **For Site Reliability Engineers**
- [**Local k3d Installation**](installation/local-k3d.md) - Kubernetes deployment guide
- [**Networking & DNS**](developers/networking-and-dns.md) - DNS and network configuration
- [**Troubleshooting Guide**](installation/local-k3d.md#troubleshooting) - Common issues and solutions

---

## 🚀 Quick Navigation

| I want to... | Go to |
|---------------|-------|
| **Understand what Fern Platform does** | [Product Overview](product/overview.md) |
| **Get started quickly (< 15 min)** | [Quick Start Guide](developers/quick-start.md) |
| **Set up local Kubernetes** | [Local k3d Setup](installation/local-k3d.md) |
| **Deploy to production** | [Local k3d Setup](installation/local-k3d.md) |
| **Configure authentication** | [OAuth Configuration](configuration/oauth.md) |
| **Understand the architecture** | [Architecture Document](ARCHITECTURE.md) |
| **Troubleshoot issues** | [Troubleshooting Guide](installation/local-k3d.md#troubleshooting) |

---

## 📖 Document Structure

```
docs/
├── product/              # Business & product information
│   ├── overview.md       # What is Fern Platform
│   ├── capabilities.md   # Features and admin functions
│   └── deployment-options.md
├── installation/         # Installation guides
│   ├── local-k3d.md      # Local Kubernetes setup
│   └── production.md     # Production deployment
├── developers/           # Developer-focused guides
│   ├── README.md         # Developer documentation hub
│   ├── quick-start.md    # 15-minute setup
│   ├── api-reference.md  # REST API reference
│   └── networking-and-dns.md
├── configuration/        # Configuration guides
│   ├── oauth.md          # OAuth setup
│   └── environment.md    # Environment variables
├── graphql-api.md       # GraphQL API documentation
├── ARCHITECTURE.md      # System architecture
├── UI_ENHANCEMENTS.md   # UI features documentation
└── rfc/                  # Technical RFCs
    └── ...
```

---

## 🔄 Migration from Old Docs

The documentation has been reorganized for clarity. Here's where to find content:

| Old Location | New Location |
|--------------|--------------|
| `README.md` (getting started) | [Quick Start Guide](developers/quick-start.md) |
| `DEPLOYMENT.md` | [Deployment Overview](../DEPLOYMENT.md) |
| `docs/k3d-deployment-guide.md` | [Local k3d Setup](installation/local-k3d.md) |
| `docs/OAuth-Setup-and-Testing.md` | [OAuth Configuration](configuration/oauth.md) |
| `docs/ARCHITECTURE.md` | [Architecture](ARCHITECTURE.md) |
| `docs/UI_ENHANCEMENTS.md` | [UI Enhancements](UI_ENHANCEMENTS.md) |

---

## 💡 Contributing to Documentation

- **Found an issue?** Open an [issue](../../issues) with the `documentation` label
- **Want to improve content?** See our [Contributing Guide](../CONTRIBUTING.md)
- **Need something not covered?** Start a [discussion](../../discussions)

---

*Last updated: June 2025*