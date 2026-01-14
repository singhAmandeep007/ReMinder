import { ReminderState, TReminder } from "types";

import { urlPrefix } from "utils";

import { TAppMockServer } from "../types";

import { resourceNotFoundResponse } from "../utils";

export function reminderRoutes(this: TAppMockServer) {
  this.get(urlPrefix("/reminders"), (schema, request) => {
    let reminders; // Collection<Instantiate<TAppRegistry, "reminder">>;
    const queryParams = request.queryParams;
    // queryParams = { groupId: "group-id-1" }

    if (queryParams) {
      if (queryParams.groupId) {
        const reminderGroup = schema.find("reminderGroup", queryParams.groupId as string);
        if (reminderGroup === null) {
          return resourceNotFoundResponse("Reminder group " + queryParams.groupId);
        }
      }
      reminders = schema.where("reminder", queryParams);
    } else {
      reminders = schema.all("reminder");
    }

    return {
      data: reminders.models.map((reminder) => reminder.getAttributes()),
    };
  });

  // NOTE: :segmentName to define a dynamic segment in the URL for a route handler
  this.delete(urlPrefix("/reminders/:id"), (schema, request) => {
    // NOTE: Access dynamic segments via request.params.segmentName
    const id = request.params.id;

    const reminder = schema.find("reminder", id);

    reminder?.destroy();

    if (reminder === null) {
      return resourceNotFoundResponse("Reminder " + id);
    }

    // NOTE: use Response class to return a custom response
    return { message: `Reminder with id ${id} deleted!` };
  });

  this.post(urlPrefix("/reminders"), (schema, request) => {
    const attrs = JSON.parse(request.requestBody);

    const defaultAttr: Partial<TReminder> = {
      isPinned: false,
      state: ReminderState.INACTIVE,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      dueDate: null,
      repeatInterval: null,
      repeatTimes: null,
      currentFocusSession: null,
      focusSessions: [],
    };

    const reminder = schema.create("reminder", Object.assign(defaultAttr, attrs));

    return {
      data: reminder.getAttributes(),
    };
  });

  this.patch(urlPrefix("/reminders/:id"), (schema, request) => {
    const id = request.params.id;
    const attrs = JSON.parse(request.requestBody);

    const reminder = schema.find("reminder", id);

    if (reminder === null) {
      return resourceNotFoundResponse("Reminder " + id);
    }

    reminder.update(attrs);

    return {
      data: reminder.getAttributes(),
    };
  });
}
