import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { HydratedDocument } from "mongoose";

export type ReminderGroupDocument = HydratedDocument<ReminderGroup>;

/**
 * Reminder group document schema
 */
@Schema({
  collection: "reminder_groups",
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
export class ReminderGroup {
  // Virtual id field
  id: string;

  @Prop({ required: true, minlength: 1 })
  name: string;

  @Prop()
  createdAt: Date;

  @Prop()
  updatedAt: Date;
}

export const ReminderGroupSchema = SchemaFactory.createForClass(ReminderGroup);

// Indexes
ReminderGroupSchema.index({ name: 1 });
ReminderGroupSchema.index({ createdAt: -1 });
