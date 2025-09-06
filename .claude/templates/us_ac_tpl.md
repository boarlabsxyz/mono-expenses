<!-- Acceptance Criteria Begin -->
## Acceptance Criteria:

### Scenario 1: [Primary Success Path]
```Gherkin
Given [Initial context]
When [Specific action]
Then [Expected result]
  And [Additional expected result]
```
#### Flow
**`tests/e2e/testdata/flows/basic/single-step.json`**:
```json [Full example of json that should be used to test the scenario]
  {
  "id": "single-step",
  "name": "Single Step Test",
  "description": "Minimal flow with one step for testing basic execution",
  "ai": {
    "provider": "openrouter",
    "apiKey": "{env.OPENROUTER_API_KEY}"
  },
  "initialStep": "only-step",
  "steps": {
    "only-step": {
      "type": "prompt",
      "prompt": "This is a test prompt",
      "model": "openrouter/auto"
    }
  }
}
```
#### Example Command
```bash  [Full example of command need to be provided]
./bin/mono-expenses execute tests/e2e/testdata/flows/integration-flows github-get-issue.json --log-level debug
```
Expected output
####
```text

```

### Scenario 2: [Error/Edge Case - if applicable]
[Same format as above]

### Scenario 3: [Alternative Path - if applicable]
[Same format as above]

*Note: Add additional scenarios as needed for comprehensive coverage*

## Success Criteria
- [Success criteria 1]
- [Success criteria 2]
- [Success criteria 3]

## Notes
- [Note 1]
- [Note 2]
- [Note 3]
<!-- Acceptance Criteria End -->
