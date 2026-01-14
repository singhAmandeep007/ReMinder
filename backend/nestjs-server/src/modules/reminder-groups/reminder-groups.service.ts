import { Injectable, Inject, NotFoundException } from "@nestjs/common";
import type { ReminderGroup, CreateReminderGroupRequest, UpdateReminderGroupRequest } from "../../generated";
import {
  IReminderGroupRepository,
  REMINDER_GROUP_REPOSITORY,
} from "../../database/interfaces/reminder-group-repository.interface";
import { IReminderRepository, REMINDER_REPOSITORY } from "../../database/interfaces/reminder-repository.interface";

/**
 * Service for managing reminder groups
 * Uses OpenAPI-generated types for consistency with API contract
 */
@Injectable()
export class ReminderGroupsService {
  constructor(
    @Inject(REMINDER_GROUP_REPOSITORY)
    private readonly reminderGroupRepository: IReminderGroupRepository,
    @Inject(REMINDER_REPOSITORY)
    private readonly reminderRepository: IReminderRepository
  ) {}

  /**
   * Get all reminder groups
   */
  async findAll(): Promise<ReminderGroup[]> {
    return this.reminderGroupRepository.findAll();
  }

  /**
   * Get a single reminder group by ID
   */
  async findOne(id: string): Promise<ReminderGroup> {
    const group = await this.reminderGroupRepository.findById(id);
    if (!group) {
      throw new NotFoundException(`Reminder group with id ${id} not found`);
    }
    return group;
  }

  /**
   * Create a new reminder group
   */
  async create(createReminderGroupDto: CreateReminderGroupRequest): Promise<ReminderGroup> {
    return this.reminderGroupRepository.create(createReminderGroupDto);
  }

  /**
   * Update an existing reminder group
   */
  async update(id: string, updateReminderGroupDto: UpdateReminderGroupRequest): Promise<ReminderGroup> {
    const existing = await this.reminderGroupRepository.findById(id);
    if (!existing) {
      throw new NotFoundException(`Reminder group with id ${id} not found`);
    }

    const updated = await this.reminderGroupRepository.update(id, updateReminderGroupDto);
    if (!updated) {
      throw new NotFoundException(`Reminder group with id ${id} not found`);
    }

    return updated;
  }

  /**
   * Delete a reminder group
   * This also deletes all reminders in the group (cascade delete)
   */
  async remove(id: string): Promise<void> {
    const existing = await this.reminderGroupRepository.findById(id);
    if (!existing) {
      throw new NotFoundException(`Reminder group with id ${id} not found`);
    }

    // Cascade delete: remove all reminders in this group
    await this.reminderRepository.deleteByGroupId(id);

    // Delete the group itself
    const deleted = await this.reminderGroupRepository.delete(id);
    if (!deleted) {
      throw new NotFoundException(`Reminder group with id ${id} not found`);
    }
  }
}
