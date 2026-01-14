/**
 * Base Repository Interface
 *
 * Defines the contract for all repository implementations.
 * Uses types from the OpenAPI-generated types.
 */

/**
 * Base repository interface for common CRUD operations
 * @template TEntity - The entity type (from OpenAPI schema)
 * @template TCreateDto - The create DTO type (from OpenAPI schema)
 * @template TUpdateDto - The update DTO type (from OpenAPI schema)
 * @template TFilter - Optional filter type for queries
 */
export interface IBaseRepository<TEntity, TCreateDto, TUpdateDto, TFilter = Record<string, unknown>> {
  /**
   * Find all entities with optional filtering
   */
  findAll(filter?: TFilter): Promise<TEntity[]>;

  /**
   * Find a single entity by ID
   */
  findById(id: string): Promise<TEntity | null>;

  /**
   * Create a new entity
   */
  create(data: TCreateDto): Promise<TEntity>;

  /**
   * Update an existing entity
   */
  update(id: string, data: Partial<TUpdateDto>): Promise<TEntity | null>;

  /**
   * Delete an entity by ID
   */
  delete(id: string): Promise<boolean>;

  /**
   * Count entities with optional filtering
   */
  count(filter?: TFilter): Promise<number>;
}
