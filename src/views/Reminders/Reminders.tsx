import { FC } from "react";

import { Panel, PanelGroup, PanelResizeHandle } from "react-resizable-panels";
import { EllipsisVertical, EllipsisIcon } from "lucide-react";

import { useBreakpoint } from "hooks";
import { BREAKPOINTS } from "config";
import { cn } from "utils";

import { RemindersList } from "./RemindersList";

import { ReminderGroupsList } from "./ReminderGroupsList";

export type TRemindersProps = Record<string, never>;

export const Reminders: FC<TRemindersProps> = () => {
  const { isBelowMd } = useBreakpoint<keyof typeof BREAKPOINTS>("md", BREAKPOINTS);

  return (
    <div
      className="mx-auto h-full max-w-screen-lg p-2 md:p-8"
      data-testid="reminders"
    >
      <div className="flex h-full max-h-[600px] flex-col rounded border-2 border-primary shadow-lg lg:h-3/4">
        <PanelGroup
          direction={isBelowMd ? "vertical" : "horizontal"}
          className="flex-1"
          autoSaveId={"reminders-panel-group"}
        >
          <Panel
            maxSize={50}
            className="flex flex-col"
          >
            <ReminderGroupsList />
          </Panel>

          <PanelResizeHandle className="group">
            <div
              className={cn(
                "flex h-full w-full items-center border-primary group-hover:bg-secondary group-hover:text-primary",
                isBelowMd ? "border-b-2" : "border-r-2"
              )}
            >
              {isBelowMd ? <EllipsisIcon className="icon mx-auto" /> : <EllipsisVertical className="icon mx-auto" />}
            </div>
          </PanelResizeHandle>

          <Panel className="flex flex-col">
            <RemindersList />
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
};
