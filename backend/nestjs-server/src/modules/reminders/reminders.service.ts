import { Injectable, Inject, NotFoundException } from "@nestjs/common";
import type { Reminder, ReminderState, CreateReminderRequest, UpdateReminderRequest } from "../../generated";
import { IReminderRepository, REMINDER_REPOSITORY } from "../../database/interfaces/reminder-repository.interface";
import {
  IReminderGroupRepository,
  REMINDER_GROUP_REPOSITORY,
} from "../../database/interfaces/reminder-group-repository.interface";

/**
 * Query parameters for listing reminders
 */
export interface QueryRemindersParams {
  groupId?: string;
  state?: ReminderState;
}

/**
 * Service for managing reminders
 * Uses OpenAPI-generated types for consistency with API contract
 */
@Injectable()
export class RemindersService {
  constructor(
    @Inject(REMINDER_REPOSITORY)
    private readonly reminderRepository: IReminderRepository,
    @Inject(REMINDER_GROUP_REPOSITORY)
    private readonly reminderGroupRepository: IReminderGroupRepository
  ) {}

  /**
   * Get all reminders with optional filtering
   */
  async findAll(query: QueryRemindersParams): Promise<Reminder[]> {
    const { groupId, state } = query;

    if (groupId && state) {
      return this.reminderRepository.findByGroupIdAndState(groupId, state);
    }

    if (groupId) {
      return this.reminderRepository.findByGroupId(groupId);
    }

    if (state) {
      return this.reminderRepository.findByState(state);
    }

    return this.reminderRepository.findAll();
  }

  /**
   * Get a single reminder by ID
   */
  async findOne(id: string): Promise<Reminder> {
    const reminder = await this.reminderRepository.findById(id);
    if (!reminder) {
      throw new NotFoundException(`Reminder with id ${id} not found`);
    }
    return reminder;
  }

  /**
   * Create a new reminder
   */
  async create(createReminderDto: CreateReminderRequest): Promise<Reminder> {
    // Validate group exists if provided
    if (createReminderDto.groupId) {
      const group = await this.reminderGroupRepository.findById(createReminderDto.groupId);
      if (!group) {
        throw new NotFoundException(`Reminder group with id ${createReminderDto.groupId} not found`);
      }
    }

    return this.reminderRepository.create(createReminderDto);
  }

  /**
   * Update an existing reminder
   */
  async update(id: string, updateReminderDto: UpdateReminderRequest): Promise<Reminder> {
    // Check if reminder exists
    const existing = await this.reminderRepository.findById(id);
    if (!existing) {
      throw new NotFoundException(`Reminder with id ${id} not found`);
    }

    // Validate group exists if changing group
    if (updateReminderDto.groupId) {
      const group = await this.reminderGroupRepository.findById(updateReminderDto.groupId);
      if (!group) {
        throw new NotFoundException(`Reminder group with id ${updateReminderDto.groupId} not found`);
      }
    }

    const updated = await this.reminderRepository.update(id, updateReminderDto);
    if (!updated) {
      throw new NotFoundException(`Reminder with id ${id} not found`);
    }

    return updated;
  }

  /**
   * Delete a reminder
   */
  async remove(id: string): Promise<void> {
    const deleted = await this.reminderRepository.delete(id);
    if (!deleted) {
      throw new NotFoundException(`Reminder with id ${id} not found`);
    }
  }
}
