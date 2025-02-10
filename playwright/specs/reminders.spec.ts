import { urlPrefix } from "utils";

import { test, expect } from "../fixtures";

import { RemindersElements } from "../pages";

test.describe("Reminders Page", () => {
  test("should have title", async ({ page }) => {
    const remindersElements = new RemindersElements(page);
    await remindersElements.goto();

    await expect(remindersElements.page).toHaveTitle(/ReMinder/);
  });

  test("should render reminders page and able to create reminder group", async ({ page, mocker }) => {
    await mocker.start();

    const remindersElements = new RemindersElements(page);
    await remindersElements.goto();

    await expect(remindersElements.root).toBeVisible();

    remindersElements.createReminderGroup("Test Group");

    await expect(remindersElements.getReminderGroupItemByText("Test Group")).toBeVisible();

    await mocker.add(
      mocker.http.post(
        urlPrefix("/reminder-groups"),
        () => {
          return mocker.HttpResponse.json(null, { status: 500 });
        },
        { once: true }
      )
    );

    remindersElements.createReminderGroup("New Test Group");

    await expect(remindersElements.page.getByText("Error creating reminder group").first()).toBeVisible();

    remindersElements.createReminderGroup("New Test Group");

    await expect(remindersElements.getReminderGroupItemByText("New Test Group")).toBeVisible();

    remindersElements.getReminderGroupMenuBtnByText("New Test Group").click();

    await expect(remindersElements.getEditReminderGroupMenuItem("New Test Group")).toHaveRole("menuitem");
  });
});
