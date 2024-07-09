import React from "react";
import { createRoot } from "react-dom/client";

import { i18n } from "modules/i18n";
import { MOCKER_TYPE } from "services/mocker";

import { App } from "./app";

import reportWebVitals from "./reportWebVitals";

import "./index.css";

async function setupApp() {
  await i18n.configure();

  const mockerType = process.env.REACT_APP_MOCKER;
  const isDemoMode = Boolean(process.env.REACT_APP_IS_DEMO_MODE);

  if (isDemoMode && !!mockerType && !window.Cypress) {
    const mocker = await import("services/mocker");

    const { setupMocker } = mocker;

    await setupMocker({ type: mockerType, shouldSeedData: true });
    // eslint-disable-next-line no-console
    console.log(`%c API is being mocked using ${mockerType}!`, "color: #bada55; font-weight: bold;");
  }

  if (window.Cypress && window.Cypress.env("REACT_APP_MOCKER") === MOCKER_TYPE.mirage) {
    import("services/mocker/mirage/proxyServer").then(({ startProxyMirageServer }) => {
      startProxyMirageServer();
    });
  }

  return Promise.resolve();
}

const root = createRoot(document.getElementById("root") as HTMLElement);

setupApp()
  .then(() => {
    root.render(
      <React.StrictMode>
        <App />
      </React.StrictMode>
    );
  })
  .catch((error) => {
    // eslint-disable-next-line no-console
    console.error("Something went wrong in setting up app", error);
  });

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
