# Security Policy

We take the security of CloudMoor seriously and appreciate responsible disclosure efforts. This document explains how to report vulnerabilities and what to expect once you do.

## Reporting a Vulnerability

Please email your report to [security@cloudmoor.dev](mailto:security@cloudmoor.dev) with the following information:

- A description of the vulnerability and its potential impact.
- Steps to reproduce (proof-of-concept code or logs are extremely helpful).
- Any mitigation ideas you may have already explored.
- Whether the issue has been disclosed elsewhere.

If you prefer, you can also use the **“Report a vulnerability”** button under GitHub’s Security tab, which routes directly to the same inbox.

> **Do not** open public issues or pull requests for security vulnerabilities. We will coordinate a disclosure timeline with you once a fix is ready.

## Response Expectations

- **Acknowledgement:** within 3 business days.
- **Initial assessment:** within 7 business days, including severity classification and remediation plan.
- **Status updates:** at least weekly until a fix is shipped.
- **Credit:** with your permission, we will credit reporters in release notes once the vulnerability is resolved.

If you have not received a response within the acknowledgement window, please feel free to follow up. We value persistence!

## Scope

This policy covers all CloudMoor source code, build scripts, configuration, and documentation in this repository as well as official distribution artefacts produced by our CI pipelines. Vulnerabilities in third-party dependencies should be reported to their maintainers, although we are happy to help coordinate if needed.

## Safe Harbor

We support research conducted in good faith. When testing, avoid actions that could disrupt production services or leak user data. As long as you follow responsible disclosure practices, we will not pursue legal action.

Thank you for keeping CloudMoor users safe.
