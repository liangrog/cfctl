cfctl
===
AWS Cloudformation DevOp tool

Design Principles & Requirements
---
1. Retain CloudFormation's independence to the tool itself
2. Dynamic, on-the-fly state management without the need to use persistent media
3. Facilitate modularity
4. Allow multiple sources of parameters
5. Unit tested all components
6. Providing three levels stack management such as module -> stack -> master stack
7. Providing parameter scoping
