import { Controller, Get, Post, Patch, Delete, Body, Param, Query, HttpCode, HttpStatus } from "@nestjs/common";
import { ApiTags, ApiOperation, ApiResponse, ApiQuery } from "@nestjs/swagger";
import type { Reminder, ReminderState, CreateReminderRequest, UpdateReminderRequest } from "../../generated";
import { RemindersService, QueryRemindersParams } from "./reminders.service";

/**
 * Controller for reminder operations
 * API contract defined in openapi/api.yaml
 */
@ApiTags("Reminders")
@Controller({ path: "reminders", version: "1" })
export class RemindersController {
  constructor(private readonly remindersService: RemindersService) {}

  /**
   * List reminders with optional filtering
   */
  @Get()
  @ApiOperation({ summary: "List reminders" })
  @ApiQuery({ name: "groupId", required: false, type: String })
  @ApiQuery({ name: "state", required: false, enum: ["INACTIVE", "INPROGRESS", "COMPLETED", "ARCHIVED"] })
  @ApiResponse({ status: 200, description: "Successful operation" })
  async findAll(@Query("groupId") groupId?: string, @Query("state") state?: ReminderState): Promise<Reminder[]> {
    const query: QueryRemindersParams = { groupId, state };
    return this.remindersService.findAll(query);
  }

  /**
   * Create a new reminder
   */
  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: "Create a reminder" })
  @ApiResponse({ status: 201, description: "Reminder created successfully" })
  @ApiResponse({ status: 400, description: "Invalid request body" })
  @ApiResponse({ status: 404, description: "Reminder group not found" })
  async create(@Body() createReminderDto: CreateReminderRequest): Promise<Reminder> {
    return this.remindersService.create(createReminderDto);
  }

  /**
   * Update an existing reminder
   */
  @Patch(":id")
  @ApiOperation({ summary: "Update a reminder" })
  @ApiResponse({ status: 200, description: "Reminder updated successfully" })
  @ApiResponse({ status: 404, description: "Reminder not found" })
  async update(@Param("id") id: string, @Body() updateReminderDto: UpdateReminderRequest): Promise<Reminder> {
    return this.remindersService.update(id, updateReminderDto);
  }

  /**
   * Delete a reminder
   */
  @Delete(":id")
  @ApiOperation({ summary: "Delete a reminder" })
  @ApiResponse({ status: 200, description: "Reminder deleted successfully" })
  @ApiResponse({ status: 404, description: "Reminder not found" })
  async remove(@Param("id") id: string): Promise<{ message: string }> {
    await this.remindersService.remove(id);
    return { message: `Reminder ${id} deleted successfully` };
  }
}
