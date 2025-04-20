import { HttpResponse, http } from "msw";

import {
  TCreateReminderGroupRequestPayload,
  TCreateReminderGroupResponsePayload,
  TDeleteReminderGroupRequestPathParams,
  TDeleteReminderGroupResponsePayload,
  TUpdateReminderGroupRequestPayload,
  TUpdateReminderGroupResponsePayload,
  TGetReminderGroupRequestPathParams,
  TGetReminderGroupResponsePayload,
  TReminderGroup,
  THTTPError,
} from "types";

import { urlPrefix } from "utils";

import { TSetupController } from "./types";

export const setupReminderGroupsController: TSetupController = ({ db }) => {
  const getReminderGroups = http.get(urlPrefix("/reminder-groups"), () => {
    return HttpResponse.json({ data: db.reminderGroup.getAll() }, { status: 200 });
  });

  const getReminderGroup = http.get<
    { id: TGetReminderGroupRequestPathParams },
    never,
    TGetReminderGroupResponsePayload | THTTPError
  >(urlPrefix("/reminder-groups/:id"), ({ params }) => {
    const reminderGroup = db.reminderGroup.findFirst({
      where: {
        id: {
          equals: params.id,
        },
      },
    });

    if (!reminderGroup) {
      return HttpResponse.json({ message: `Reminder Group with id ${params.id} not found!` }, { status: 404 });
    }

    return HttpResponse.json({ data: reminderGroup }, { status: 200 });
  });

  const createReminderGroup = http.post<
    never,
    TCreateReminderGroupRequestPayload,
    TCreateReminderGroupResponsePayload | THTTPError
  >(urlPrefix("/reminder-groups"), async ({ request }) => {
    const requestBody = await request.json();

    const reminderGroup = db.reminderGroup.create({
      ...requestBody,
    });

    return HttpResponse.json({ data: reminderGroup }, { status: 201 });
  });

  const deleteReminderGroup = http.delete<
    { id: TDeleteReminderGroupRequestPathParams },
    never,
    TDeleteReminderGroupResponsePayload | THTTPError
  >(urlPrefix("/reminder-groups/:id"), ({ params }) => {
    const reminderGroup = db.reminderGroup.findFirst({
      where: {
        id: {
          equals: params.id,
        },
      },
    });

    if (!reminderGroup) {
      return HttpResponse.json({ message: `Reminder Group with id ${params.id} not found!` }, { status: 404 });
    }

    db.reminder.deleteMany({
      where: {
        group: {
          id: {
            equals: params.id,
          },
        },
      },
    });

    db.reminderGroup.delete({
      where: {
        id: {
          equals: params.id,
        },
      },
    });

    return HttpResponse.json({ message: `Reminder Group with id ${params.id} deleted!` }, { status: 200 });
  });

  const updateReminderGroup = http.patch<
    { id: TUpdateReminderGroupRequestPayload["id"] },
    TUpdateReminderGroupRequestPayload,
    TUpdateReminderGroupResponsePayload | THTTPError
  >(urlPrefix("/reminder-groups/:id"), async ({ request, params }) => {
    const requestBody = await request.json();

    const reminderGroup = db.reminderGroup.findFirst({
      where: {
        id: {
          equals: params.id,
        },
      },
    });

    if (!reminderGroup) {
      return HttpResponse.json({ message: `Reminder Group with id ${params.id} not found!` }, { status: 404 });
    }

    const updatedReminder = db.reminderGroup.update({
      where: {
        id: {
          equals: params.id,
        },
      },
      data: {
        ...requestBody,
      },
    }) as TReminderGroup;

    return HttpResponse.json({ data: updatedReminder }, { status: 200 });
  });

  return [getReminderGroups, getReminderGroup, createReminderGroup, deleteReminderGroup, updateReminderGroup];
};
