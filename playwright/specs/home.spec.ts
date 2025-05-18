import { test, expect } from "@playwright/test";

import { HomeElements } from "../pages";

test.describe("Home Page", () => {
  test("should have title", async ({ page }) => {
    const homeElements = new HomeElements(page);
    await homeElements.goto();

    await expect(homeElements.page).toHaveTitle(/ReMinder/);
  });

  test("should render home page content and navigate to reminders page", async ({ page }) => {
    const homeElements = new HomeElements(page);

    await homeElements.goto();

    await expect(homeElements.root.getByRole("heading", { name: "Full-Stack Reminder Application" })).toBeVisible();

    await expect(homeElements.root.getByRole("link", { name: "Explore Demo" })).toBeVisible();

    await expect(homeElements.root.getByRole("link", { name: "Check Build Process" })).toBeVisible();

    await expect(homeElements.footer).toContainText(`Copyright © ${new Date().getFullYear()}Amandeep Singh`);

    await expect(homeElements.footer.getByRole("link", { name: "Amandeep Singh" })).toHaveAttribute(
      "href",
      "https://github.com/singhAmandeep007"
    );

    await homeElements.header.getByRole("link", { name: "Explore Demo" }).click();

    await expect(homeElements.root).toBeHidden();
  });
});
