## Mannual mocks

1. Docs:
   1. https://jestjs.io/docs/manual-mocks
1. Issues:
   1. https://github.com/facebook/create-react-app/issues/7539

## MSW

1. [Setup](./src/services/mocker/msw)
   1. [DB](./src/services/mocker/msw/db.ts)
   2. [Handlers](./src/services/mocker/msw/handlers.ts)
   3. [Server](./src/services/mocker/msw/server.ts)

### Setting up MSW with Jest(JSDOM)

1. Docs:
   1. https://github.com/mswjs/examples/blob/main/examples/with-jest-jsdom/README.md
   2. https://mswjs.io/docs/migrations/1.x-to-2.x/#requestresponsetextencoder-is-not-defined-jest
2. Implementation:
   1. [testServer.ts](./src/tests/utils/testServer.ts)

## RTL

1. Wrapper for RTL -
   1. [Component](./src/tests/utils/render.tsx)
   2. [Hooks](./src/tests/utils/renderHook.tsx)
2. Notes:
   1. Print whole JSDOM object - `screen.debug(result.container, Infinity);`
   2. When to use `act` - https://github.com/threepointone/react-act-examples/blob/master/sync.md

## Mirage

1. [Setup](./src/services/mocker/mirage)
   1. [Routes](./src/services/mocker/mirage/routes/index.ts)
   2. [Proxy Server](./src/services/mocker/mirage/proxyServer.ts)
   3. [Server](./src/services/mocker/mirage/server.ts)
   4. [Models](./src/services/mocker/mirage/models/index.ts)
   5. [Factories](./src/services/mocker/mirage/factories/index.ts)
2. Notes:
   1. Print all entities in db - `console.log(testMirageServer.db.dump());`
   2. [Track requests for assertions](https://miragejs.com/docs/testing/assertions/#asserting-against-handled-requests-and-responses)

### Setting up Mirage with Jest(JSDOM)

1. Docs:
   1. https://miragejs.com/docs/testing/integration-and-unit-tests/
   2. https://miragejs.com/tutorial/part-9/
2. Implementation:
   1. [testServer.ts](./src/tests/utils/testServer.ts)


## MSW + Cypress

1. https://github.com/mswjs/msw/issues/1560
2. [Implementation](./cypress/support/msw.ts)

## Mirage + Cypress

1. https://miragejs.com/quickstarts/cypress/
2. [Implementation](./cypress/support/mirage.ts)

## Popper issue

1. https://github.com/floating-ui/react-popper/issues/350

## Conventions followed

### test ids

1. testId - `data-testid`
2. kebab case - `long-test-id`
3. single test id as prop - `testId: 'id'`
4. multiple test id as prop - `testIds: { id1: 'id-1', id2: 'id-2' }`

## E2E

### Cypress

[Config](./cypress.config.ts)

1. headless
   - `npm run cy:test -- --browser chrome` - Run tests in chrome browser in headless mode
   - `npm run cy:test -- --spec "cypress/specs/home.cy.ts"` - Run specific test file in headless mode
2. headed
   - `npm run cy:test:open` - Run tests in headed mode, closes on completion
   - `npm run cy:test:open -- --spec "cypress/specs/home.cy.ts"` - Run specific test file in headed mode
3. launchpad
   1. first run `npm run cy:server:start` - Starts dev server with cypress env variables
   2. then run `npm run cy:open` in a new terminal - Opens cypress launchpad (need to be closed manually).
4. change `REACT_APP_MOCKER` in `.cypress.env` to switch between `msw` and `mirage`. Note: restart cypress after changing the value.

Reports:

1. Cypress Test Report
   1. [HTML Report](./reports/cypress/result/html)
   2. [JSON Report](./reports/cypress/result/json)
   3. [JUnit Report](./reports/cypress/result/xml)
2. [Coverage Report](./reports/cypress/coverage)

### Playwright

[Config](./playwright.config.ts)

1. headless
   - `npm run pw:test` - Run tests in headless mode
   - `npm run pw:test specs/home.spec.ts` - Run specific test file in headless mode
2. ui mode
   - `npm run pw:open` - Stars dev server and then opens launchpad.
3. Only msw as a mocker is supported in playwright.


Reports:

1. Playwright Test Report
   1. [HTML Report](./reports/playwright/result/html)
   2. [JSON Report](./reports/playwright/result/json)
   3. [JUnit Report](./reports/playwright/result/xml)
