import { Factory } from "miragejs";
// eslint-disable-next-line import/no-unresolved
import { FactoryDefinition } from "miragejs/-types";

import { TReminder, ReminderState } from "types";

import { TAppModels } from "../types";

export const reminder: FactoryDefinition<TAppModels["reminder"]> = Factory.extend<TReminder>({
  id(n) {
    return `reminder-id-${n}`;
  },
  title(n) {
    return `Reminder ${n}`;
  },
  isPinned(n) {
    if (n % 2 === 0) {
      return true;
    }
    return false;
  },
  state: ReminderState.INACTIVE,
  createdAt() {
    return new Date().toISOString();
  },
  updatedAt() {
    return new Date().toISOString();
  },
  dueDate() {
    return null;
  },
  repeatInterval() {
    return null;
  },
  repeatTimes() {
    return null;
  },
  group() {
    return null;
  },
  focusSessions() {
    return [];
  },
  currentFocusSession() {
    return null;
  },
});
