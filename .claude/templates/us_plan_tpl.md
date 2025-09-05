<!-- Implementation Plan Begin -->
## Implementation Plan

### Overview of Changes
[Brief description of what this implementation will achieve and the main approach]

### Technology Stack
- Framework: [Framework choice with version]
- Build Tool: [Build tool and configuration]
- Language: [Language and version]
- External API: [Any external services]

### Complexity
- Complexity Level: 3
- Simple Enhancement

### Technology Validation Checkpoints
- [x] Required dependencies identified and installed ([list if new])
- [x] Build configuration validated ([build tool] works)
- [x] Hello world verification completed ([POC description])
- [x] Test build passes successfully ([command used])

### Files to Modify/Create

#### 1. **[path/to/file1.ext]** ([NEW/MODIFY])
- [Purpose of this file]
- [Key functionality to implement]
- [Interfaces/types to define]

#### 2. **[path/to/file2.ext]** ([NEW/MODIFY])
- [Purpose of this file]
- [Key functionality to implement]
- [Integration points]

#### 3. **[path/to/file3.ext]** ([NEW/MODIFY])
- [Purpose of this file]
- [Changes needed]
- [Dependencies on other files]

### Implementation Steps

#### Step 1: [Write or update e2e test]
1.1. **`tests/e2e/[path_to_test]_test.go`**
   - [Test 1] - [what to verify with code example 1]
   - [Test 2] - [what to verify with code example 2]
   - [Test 3] - [what to verify  with code example 3]

1.2. **`tests/e2e/[path_to_test]_test.go`**
   - [Test 1] - [what to verify with code example 1]
   - [Test 2] - [what to verify with code example 2]
   - [Test 3] - [what to verify  with code example 3]
1.3. Run e2e test
   - Expected result: Updated/created tests should fail

#### Step 2: [Step Name]
1. [Subtask 2.1 - specific action]
2. [Subtask 2.2 - specific action]
3. [Subtask 2.3 - specific action]
4. [Subtask 2.4 - specific action]

#### Step 3: [Step Name]
1. [Subtask 3.1 - specific action]
2. [Subtask 3.2 - specific action]
3. [Subtask 3.3 - specific action]
4. [Subtask 3.4 - specific action]
5. [Subtask 3.5 - specific action]

#### Step 4: [Step Name]
1. [Subtask 4.1 - specific action]
2. [Subtask 4.2 - specific action]
3. [Subtask 4.3 - specific action]

#### Step 5: [Step Name]
1. [Subtask 5.1 - specific action]
2. [Subtask 5.2 - specific action]
3. [Subtask 5.3 - specific action]
4. [Subtask 5.4 - specific action]


#### Final Step [Number]: Validation (this must be last step, so number should be follow the order)
1. Confirm all Technology Validation Checkpoints are completed
2. Build the project: `make build`
3. Format the project: `gofumpt -l -w .`
4. Lint the project: `golangci-lint run --fix`
5. Run unit tests: `make test-unit`, fix issues
6. Run e2e tests: `make test-e2e`, fix issues

### Potential Challenges & Mitigations

1. **[Challenge Name]**
   - **Challenge**: [Description of the challenge]
   - **Mitigation**: [How to handle this challenge]

2. **[Challenge Name]**
   - **Challenge**: [Description of the challenge]
   - **Mitigation**: [How to handle this challenge]

3. **[Challenge Name]**
   - **Challenge**: [Description of the challenge]
   - **Mitigation**: [How to handle this challenge]

4. **[Challenge Name]**
   - **Challenge**: [Description of the challenge]
   - **Mitigation**: [How to handle this challenge]

5. **[Challenge Name]**
   - **Challenge**: [Description of the challenge]
   - **Mitigation**: [How to handle this challenge]

### Contracts, Schemas and Interface Updates

#### New Interface: [InterfaceName]
```go
type [InterfaceName] interface {
    [Method1]([params]) ([returns])
    [Method2]([params]) ([returns])
}
```

#### Updated [StructName] Constructor
```go
func New[StructName]([params]) *[StructName]
```

#### Updated [StructName] Fields
```go
type [StructName] struct {
    [existingField1] [type]
    [existingField2] [type]
    [newField]       [type]  // NEW field
}
```

#### No Changes Expected
- [TestName1] - [Reason why no changes needed]
- [TestName2] - [Reason why no changes needed]

### Task Checklist

#### Preparation
- [x] [Completed preparation task]
- [x] [Completed preparation task]
- [x] [Completed preparation task]
- [x] [Completed preparation task]

### Dependencies
- [Dependency 1 - with explanation]
- [Dependency 2 - with explanation]
- [Dependency 3 - with explanation]

## Planning Summary

### Key Decisions Made
1. [Decision 1 with rationale]
2. [Decision 2 with rationale]
3. [Decision 3 with rationale]
4. [Decision 4 with rationale]
5. [Decision 5 with rationale]
