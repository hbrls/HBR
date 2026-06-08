---
name: template-structure
metadata:
    version: 0.0.2
---

# Code Structure

## Project Layout

```
{{project-name}}/
├── src/                    # Source code
│   ├── {{module1}}/        # {{description}}
│   ├── {{module2}}/        # {{description}}
│   └── {{module3}}/        # {{description}}
├── tests/                  # Test files
│   ├── unit/               # Unit tests
│   ├── integration/        # Integration tests
│   └── e2e/                # End-to-end tests
├── docs/                   # Documentation
├── scripts/                # Build/deploy scripts
└── {{other}}/              # {{description}}
```

## Module Organization

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Files | {{convention}} | {{example}} |
| Directories | {{convention}} | {{example}} |
| Classes | {{convention}} | {{example}} |
| Functions | {{convention}} | {{example}} |
| Constants | {{convention}} | {{example}} |
| Interfaces | {{convention}} | {{example}} |

### Import Order

1. Standard library imports
2. Third-party imports
3. Local application imports

```{{language}}
// Example
import {{standard library}}
import {{third party}}
import {{local module}}
```

## Code Patterns

### Component Structure

<!-- Standard structure for components/modules -->

```{{language}}
{{component template}}
```

### Error Handling

<!-- Standard error handling pattern -->

```{{language}}
{{error handling template}}
```

### Logging

<!-- Standard logging pattern -->

```{{language}}
{{logging template}}
```

## Testing Conventions

### Test File Location

- Unit tests: `tests/unit/{{module}}/`
- Integration tests: `tests/integration/`
- E2E tests: `tests/e2e/`

### Test Naming

```
test_{{function_name}}_{{scenario}}_{{expected_result}}
```

### Test Structure

```{{language}}
describe('{{Component}}', () => {
  describe('{{method}}', () => {
    it('should {{expected behavior}} when {{condition}}', () => {
      // Arrange
      // Act
      // Assert
    });
  });
});
```

## Documentation Standards

### Code Comments

- Use comments to explain WHY, not WHAT
- Document public APIs with docstrings
- Keep comments up-to-date with code changes

### README Requirements

Each module should have a README with:

- Purpose description
- Quick start example
- API reference (or link)
- Dependencies

## Dependency Management

**涉及依赖变更一律转人工处理**

## Performance Guidelines

### Code-Level

- Avoid N+1 queries
- Use lazy loading for large data
- Cache expensive computations
- Profile before optimizing

### Bundle Size (Frontend)

- Maximum bundle size: {{size}}
- Use code splitting for routes
- Lazy load heavy components

## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
