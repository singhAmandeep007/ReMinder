import { Injectable } from "@nestjs/common";
import { InjectModel } from "@nestjs/mongoose";
import { Model } from "mongoose";
import type {
  ReminderGroup as IReminderGroup,
  CreateReminderGroupRequest,
  UpdateReminderGroupRequest,
} from "../../../generated";
import type {
  IReminderGroupRepository,
  IReminderGroupFilter,
} from "../../interfaces/reminder-group-repository.interface";
import { ReminderGroup, ReminderGroupDocument } from "../schemas/reminder-group.schema";

/**
 * MongoDB implementation of the ReminderGroup repository
 * Returns OpenAPI-compliant ReminderGroup types
 */
@Injectable()
export class MongoReminderGroupRepository implements IReminderGroupRepository {
  constructor(
    @InjectModel(ReminderGroup.name)
    private readonly reminderGroupModel: Model<ReminderGroupDocument>
  ) {}

  async findAll(filter?: IReminderGroupFilter): Promise<IReminderGroup[]> {
    const mongoFilter = this.buildFilter(filter);
    const groups = await this.reminderGroupModel.find(mongoFilter).sort({ createdAt: -1 }).exec();
    return groups.map((doc) => doc.toJSON() as unknown as IReminderGroup);
  }

  async findById(id: string): Promise<IReminderGroup | null> {
    try {
      const group = await this.reminderGroupModel.findById(id).exec();
      return group ? (group.toJSON() as unknown as IReminderGroup) : null;
    } catch {
      return null;
    }
  }

  async findByName(name: string): Promise<IReminderGroup | null> {
    const group = await this.reminderGroupModel.findOne({ name }).exec();
    return group ? (group.toJSON() as unknown as IReminderGroup) : null;
  }

  async create(data: CreateReminderGroupRequest): Promise<IReminderGroup> {
    const group = new this.reminderGroupModel(data);
    const saved = await group.save();
    return saved.toJSON() as unknown as IReminderGroup;
  }

  async update(id: string, data: Partial<UpdateReminderGroupRequest>): Promise<IReminderGroup | null> {
    try {
      const group = await this.reminderGroupModel.findByIdAndUpdate(id, { $set: data }, { new: true }).exec();
      return group ? (group.toJSON() as unknown as IReminderGroup) : null;
    } catch {
      return null;
    }
  }

  async delete(id: string): Promise<boolean> {
    try {
      const result = await this.reminderGroupModel.deleteOne({ _id: id }).exec();
      return result.deletedCount > 0;
    } catch {
      return false;
    }
  }

  async count(filter?: IReminderGroupFilter): Promise<number> {
    const mongoFilter = this.buildFilter(filter);
    return this.reminderGroupModel.countDocuments(mongoFilter).exec();
  }

  /**
   * Build MongoDB filter from filter interface
   */
  private buildFilter(filter?: IReminderGroupFilter): Record<string, unknown> {
    if (!filter) return {};

    const mongoFilter: Record<string, unknown> = {};

    if (filter.name) mongoFilter.name = filter.name;

    return mongoFilter;
  }
}
