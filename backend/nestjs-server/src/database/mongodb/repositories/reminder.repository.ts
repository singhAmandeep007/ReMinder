import { Injectable } from "@nestjs/common";
import { InjectModel } from "@nestjs/mongoose";
import { Model } from "mongoose";
import type {
  Reminder as IReminder,
  ReminderState,
  CreateReminderRequest,
  UpdateReminderRequest,
} from "../../../generated";
import type { IReminderRepository, IReminderFilter } from "../../interfaces/reminder-repository.interface";
import { Reminder, ReminderDocument } from "../schemas/reminder.schema";

/**
 * MongoDB implementation of the Reminder repository
 * Returns OpenAPI-compliant Reminder types
 */
@Injectable()
export class MongoReminderRepository implements IReminderRepository {
  constructor(
    @InjectModel(Reminder.name)
    private readonly reminderModel: Model<ReminderDocument>
  ) {}

  async findAll(filter?: IReminderFilter): Promise<IReminder[]> {
    const mongoFilter = this.buildFilter(filter);
    const reminders = await this.reminderModel.find(mongoFilter).sort({ isPinned: -1, createdAt: -1 }).exec();
    return reminders.map((doc) => doc.toJSON() as unknown as IReminder);
  }

  async findById(id: string): Promise<IReminder | null> {
    try {
      const reminder = await this.reminderModel.findById(id).exec();
      return reminder ? (reminder.toJSON() as unknown as IReminder) : null;
    } catch {
      return null;
    }
  }

  async create(data: CreateReminderRequest): Promise<IReminder> {
    const reminder = new this.reminderModel({
      ...data,
      state: "INACTIVE" as ReminderState,
      isPinned: false,
      focusSessions: [],
    });
    const saved = await reminder.save();
    return saved.toJSON() as unknown as IReminder;
  }

  async update(id: string, data: Partial<UpdateReminderRequest>): Promise<IReminder | null> {
    try {
      const reminder = await this.reminderModel.findByIdAndUpdate(id, { $set: data }, { new: true }).exec();
      return reminder ? (reminder.toJSON() as unknown as IReminder) : null;
    } catch {
      return null;
    }
  }

  async delete(id: string): Promise<boolean> {
    try {
      const result = await this.reminderModel.deleteOne({ _id: id }).exec();
      return result.deletedCount > 0;
    } catch {
      return false;
    }
  }

  async count(filter?: IReminderFilter): Promise<number> {
    const mongoFilter = this.buildFilter(filter);
    return this.reminderModel.countDocuments(mongoFilter).exec();
  }

  async findByGroupId(groupId: string): Promise<IReminder[]> {
    const reminders = await this.reminderModel.find({ groupId }).sort({ isPinned: -1, createdAt: -1 }).exec();
    return reminders.map((doc) => doc.toJSON() as unknown as IReminder);
  }

  async findByState(state: ReminderState): Promise<IReminder[]> {
    const reminders = await this.reminderModel.find({ state }).sort({ isPinned: -1, createdAt: -1 }).exec();
    return reminders.map((doc) => doc.toJSON() as unknown as IReminder);
  }

  async findByGroupIdAndState(groupId: string, state: ReminderState): Promise<IReminder[]> {
    const reminders = await this.reminderModel.find({ groupId, state }).sort({ isPinned: -1, createdAt: -1 }).exec();
    return reminders.map((doc) => doc.toJSON() as unknown as IReminder);
  }

  async deleteByGroupId(groupId: string): Promise<number> {
    const result = await this.reminderModel.deleteMany({ groupId }).exec();
    return result.deletedCount;
  }

  /**
   * Build MongoDB filter from filter interface
   */
  private buildFilter(filter?: IReminderFilter): Record<string, unknown> {
    if (!filter) return {};

    const mongoFilter: Record<string, unknown> = {};

    if (filter.state) mongoFilter.state = filter.state;
    if (filter.isPinned !== undefined) mongoFilter.isPinned = filter.isPinned;
    if (filter.groupId) mongoFilter.groupId = filter.groupId;

    return mongoFilter;
  }
}
