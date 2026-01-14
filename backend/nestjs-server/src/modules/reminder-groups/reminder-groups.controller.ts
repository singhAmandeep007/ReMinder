import { Controller, Get, Post, Patch, Delete, Body, Param, HttpCode, HttpStatus } from "@nestjs/common";
import { ApiTags, ApiOperation, ApiResponse } from "@nestjs/swagger";
import type { ReminderGroup, CreateReminderGroupRequest, UpdateReminderGroupRequest } from "../../generated";
import { ReminderGroupsService } from "./reminder-groups.service";

/**
 * Controller for reminder group operations
 * API contract defined in openapi/api.yaml
 */
@ApiTags("Reminder Groups")
@Controller({ path: "reminder-groups", version: "1" })
export class ReminderGroupsController {
  constructor(private readonly reminderGroupsService: ReminderGroupsService) {}

  /**
   * List all reminder groups
   */
  @Get()
  @ApiOperation({ summary: "List reminder groups" })
  @ApiResponse({ status: 200, description: "Successful operation" })
  async findAll(): Promise<ReminderGroup[]> {
    return this.reminderGroupsService.findAll();
  }

  /**
   * Get a single reminder group by ID
   */
  @Get(":id")
  @ApiOperation({ summary: "Get a reminder group" })
  @ApiResponse({ status: 200, description: "Successful operation" })
  @ApiResponse({ status: 404, description: "Reminder group not found" })
  async findOne(@Param("id") id: string): Promise<ReminderGroup> {
    return this.reminderGroupsService.findOne(id);
  }

  /**
   * Create a new reminder group
   */
  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiOperation({ summary: "Create a reminder group" })
  @ApiResponse({ status: 201, description: "Reminder group created successfully" })
  @ApiResponse({ status: 400, description: "Invalid request body" })
  async create(@Body() createReminderGroupDto: CreateReminderGroupRequest): Promise<ReminderGroup> {
    return this.reminderGroupsService.create(createReminderGroupDto);
  }

  /**
   * Update an existing reminder group
   */
  @Patch(":id")
  @ApiOperation({ summary: "Update a reminder group" })
  @ApiResponse({ status: 200, description: "Reminder group updated successfully" })
  @ApiResponse({ status: 404, description: "Reminder group not found" })
  async update(
    @Param("id") id: string,
    @Body() updateReminderGroupDto: UpdateReminderGroupRequest
  ): Promise<ReminderGroup> {
    return this.reminderGroupsService.update(id, updateReminderGroupDto);
  }

  /**
   * Delete a reminder group
   */
  @Delete(":id")
  @ApiOperation({ summary: "Delete a reminder group" })
  @ApiResponse({ status: 200, description: "Reminder group deleted successfully" })
  @ApiResponse({ status: 404, description: "Reminder group not found" })
  async remove(@Param("id") id: string): Promise<{ message: string }> {
    await this.reminderGroupsService.remove(id);
    return { message: `Reminder group ${id} deleted successfully` };
  }
}
