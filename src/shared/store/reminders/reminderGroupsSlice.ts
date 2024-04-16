import {
  TReminderGroup,
  TGetReminderGroupsResponsePayload,
  TGetReminderGroupRequestPayload,
  TGetReminderGroupResponsePayload,
  TCreateReminderGroupRequestPayload,
  TCreateReminderGroupResponsePayload,
  TUpdateReminderGroupRequestPayload,
  TUpdateReminderGroupResponsePayload,
  TDeleteReminderGroupRequestPayload,
  TDeleteReminderGroupResponsePayload,
} from "types";

import { apiSlice } from "../apiSlice";

export const reminderGroupsApiSlice = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    getReminderGroups: builder.query<TReminderGroup[], void>({
      query: () => "/reminder-groups",
      transformResponse: (response: TGetReminderGroupsResponsePayload) => response.data,
      providesTags: (result) =>
        result
          ? // successful query
            [
              ...result.map(({ id }) => ({ type: "ReminderGroups", id }) as const),
              { type: "ReminderGroups", id: "LIST" },
            ]
          : // an error occurred,
            [{ type: "ReminderGroups", id: "LIST" }],
    }),
    getReminderGroup: builder.query<TReminderGroup, TGetReminderGroupRequestPayload>({
      query: (id) => `/reminder-groups/${id}`,
      transformResponse: (response: TGetReminderGroupResponsePayload) => response.data,
      providesTags: (result, error, id) => [{ type: "ReminderGroups", id }],
    }),
    deleteReminderGroup: builder.mutation<TDeleteReminderGroupResponsePayload, TDeleteReminderGroupRequestPayload>({
      query: (id) => ({
        url: `/reminder-groups/${id}`,
        method: "DELETE",
      }),
      // FIX: invalidates getReminders query having deleted id as argument. No better way to prevent auto refetching of it
      invalidatesTags: (result, error, id) => [
        { type: "ReminderGroups", id: "LIST" },
        { type: "Reminders", id: "LIST" },
      ],
    }),
    // rr
    createReminderGroup: builder.mutation<TReminderGroup, TCreateReminderGroupRequestPayload>({
      query: (body) => ({
        url: "/reminder-groups",
        method: "POST",
        body,
      }),
      transformResponse: (response: TCreateReminderGroupResponsePayload) => response.data,
      // invalidatesTags: [{ type: "ReminderGroups", id: "LIST" }],
    }),
    updateReminderGroup: builder.mutation<TReminderGroup, TUpdateReminderGroupRequestPayload>({
      query: ({ id, ...body }) => ({
        url: `/reminder-groups/${id}`,
        method: "PATCH",
        body,
      }),
      transformResponse: (response: TUpdateReminderGroupResponsePayload) => response.data,
      // invalidatesTags: (response, error, { id }) => [
      //   { type: "ReminderGroups", id },
      //   { type: "Reminders", id: "LIST" },
      // ],
    }),
  }),
});

export const selectReminderGroups = () => reminderGroupsApiSlice.endpoints.getReminderGroups.select();
