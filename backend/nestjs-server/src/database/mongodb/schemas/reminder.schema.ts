import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { HydratedDocument, Types } from "mongoose";

/**
 * Reminder state enum
 */
export enum ReminderState {
  INACTIVE = "INACTIVE",
  INPROGRESS = "INPROGRESS",
  COMPLETED = "COMPLETED",
  ARCHIVED = "ARCHIVED",
}

export type ReminderDocument = HydratedDocument<Reminder>;

/**
 * Focus session subdocument schema
 */
@Schema({ _id: false })
export class ReminderFocusSession {
  @Prop({ required: true, type: Date })
  startTime: Date;

  @Prop({ required: true, type: Date })
  endTime: Date;
}

export const ReminderFocusSessionSchema = SchemaFactory.createForClass(ReminderFocusSession);

/**
 * Reminder document schema
 */
@Schema({
  collection: "reminders",
  timestamps: true,
  toJSON: {
    virtuals: true,
    transform: (_, ret: Record<string, unknown>) => {
      ret.id = String(ret._id);
      delete ret._id;
      delete ret.__v;
      return ret;
    },
  },
})
export class Reminder {
  // Virtual id field
  id: string;

  @Prop({ required: true, minlength: 1 })
  title: string;

  @Prop({
    required: true,
    enum: Object.values(ReminderState),
    default: ReminderState.INACTIVE,
  })
  state: ReminderState;

  @Prop({ required: true, default: false })
  isPinned: boolean;

  @Prop({ type: Types.ObjectId, ref: "ReminderGroup", default: null })
  groupId: Types.ObjectId | null;

  @Prop({ type: Date, default: null })
  dueDate: Date | null;

  @Prop({ type: Number, default: null, min: 1 })
  repeatTimes: number | null;

  @Prop({ type: Number, default: null, min: 0 })
  repeatInterval: number | null;

  @Prop({ type: ReminderFocusSessionSchema, default: null })
  currentFocusSession: ReminderFocusSession | null;

  @Prop({ type: [ReminderFocusSessionSchema], default: [] })
  focusSessions: ReminderFocusSession[];

  @Prop()
  createdAt: Date;

  @Prop()
  updatedAt: Date;
}

export const ReminderSchema = SchemaFactory.createForClass(Reminder);

// Indexes for performance
ReminderSchema.index({ state: 1 });
ReminderSchema.index({ groupId: 1 }, { sparse: true });
ReminderSchema.index({ createdAt: -1 });
ReminderSchema.index({ isPinned: -1, createdAt: -1 });
ReminderSchema.index({ title: "text" });
