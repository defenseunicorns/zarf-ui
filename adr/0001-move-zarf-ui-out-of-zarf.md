# 1. Move Zarf UI out of Zarf

Date: 2023-09-30

## Status

Accepted

## Context

One of Zarf's primary personas has been (and will continue to be) Ashton, a junior DevOps SRE with minimal experience managing services in production.  To serve this persona Zarf has focused (and will continue to focus) on User Experience for creating and deploying packages along with managing deployments on Day 2.  One experiment that came from this was the Zarf Web UI that embedded a Web User Interface into Zarf and allowed a user to execute deployments from their browser instead of their CLI.  As this experiment grew however it was realized that Zarf's binary may not be the best home for this UI for the following reasons:

1. The Zarf team itself has a set amount of bandwidth and to take the UI where we want it to go would require more resources.
2. Running the Zarf UI from the Zarf binary directly makes it difficult to deploy a long-lived version of the UI that could provide more advanced capabilities.

## Decision

To that end, the Zarf UI will be migrated to become a capability of the Unicorn Delivery Service so that it can receive dedicated support and a better architectural paradigm.  Given this team is new and is just standing up, the Zarf UI will be temporarily moved to another repository (https://github.com/defenseunicorns/zarf-ui) where maintenance issues and smaller PRs can continue to be worked and released separate from the main Zarf project allowing them to focus on a better CLI experience.

## Consequences

This will mean a slowing of updates in the near term for the Zarf UI as the temporary project will only be receiving maintenance or community updates but it will help pave the way for a better future for the UI beyond just Zarf itself.
