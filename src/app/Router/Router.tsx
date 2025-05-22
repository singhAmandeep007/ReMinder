import { FC } from "react";

import { RouterProvider, createBrowserRouter } from "react-router";

import { ErrorFallback } from "views/ErrorFallback";

import { NotFound } from "views/NotFound";

import { BaseLayout } from "views/Layout";

import { Home } from "views/Home/index";

import { HOME_ROUTE_BY_PATH, ROUTE_BY_PATH, BASENAME } from "./consts";

const router = createBrowserRouter(
  [
    {
      element: <BaseLayout />,
      errorElement: <ErrorFallback />,
      hydrateFallbackElement: <div>Loading...</div>,
      children: [
        {
          path: ROUTE_BY_PATH.home,
          index: true,
          element: <Home />,
        },
        {
          path: HOME_ROUTE_BY_PATH.reminders,
          lazy: async () => {
            const { Reminders } = await import("views/Reminders");
            return {
              Component: Reminders,
            };
          },
        },
        {
          path: ROUTE_BY_PATH.notFound,
          element: <NotFound />,
        },
      ],
    },
  ],
  {
    basename: BASENAME,
  }
);

export type TRouterProps = Record<string, never>;

export const Router: FC<TRouterProps> = () => <RouterProvider router={router} />;
