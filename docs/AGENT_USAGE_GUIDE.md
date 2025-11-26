# LiveTemplate Agent Usage Guide

This guide explains how to use the LiveTemplate assistant agent in Claude Code to create and develop applications.

## Getting Started (Recommended Workflow)

The recommended approach is to let the agent handle the entire project creation from start to finish. This allows the agent to:
- Understand the complete context and architecture
- Plan before coding
- Verify each step (tests pass, migrations work)
- Make consistent decisions throughout

```bash
# Open Claude Code in any directory
claude
```

Then describe what you want to build:
```
You: Create a new LiveTemplate blog app with posts
```

The agent will:
1. Plan the project structure and architecture
2. Run `lvt new` to create the project
3. Generate a posts resource with appropriate fields
4. Create and run migrations
5. Start the server
6. Open your browser to the running app

**Why this approach?** The agent sees the entire creation process, understands decisions made at each step, and can verify everything works before moving forward. This matches industry best practices for AI-assisted development.

---

## Alternative: Manual Project Creation

If you prefer more hands-on control, you can create the project manually first:

```bash
# Create a new LiveTemplate app
lvt new myapp

# Navigate to the project
cd myapp

# Open in Claude Code CLI
claude
```

Then ask the agent to add features:
```
You: Add a posts resource with title and content
```

The agent will generate the posts resource and run migrations. You can then start the server with `go run cmd/myapp/main.go`.

## Example Workflows

### Complete Project Creation (Recommended)

Start with a clear description of what you want:

```
You: I want to build a blog with posts that have titles, content,
     and publication dates. Include categories and tags.
```

The agent will:
1. **Plan**: Analyze requirements and propose architecture
2. **Create**: Run `lvt new blog` to initialize project
3. **Generate**: Create posts, categories, and tags resources with proper relationships
4. **Migrate**: Create and apply database migrations
5. **Verify**: Run tests and check migrations applied correctly
6. **Launch**: Start the development server

You get a fully working blog application, ready to customize.

### Full Stack Application

For more complex requirements:

```
You: Create a task management system with:
- User authentication (email/password)
- Projects that belong to users
- Tasks within projects with due dates and priorities
- Real-time updates when tasks change
```

The agent will:
1. **Plan**: Design the data model and relationships
2. **Setup**: Create project with multi kit (for multiple resources)
3. **Authenticate**: Generate complete auth system with sessions
4. **Resources**: Create projects and tasks with foreign keys
5. **Real-time**: Add WebSocket support for live updates
6. **Test**: Verify authentication flow and CRUD operations
7. **Launch**: Start server with all features working

### Incremental Development

You can also build features incrementally:

```
You: Create a simple blog app

[Agent creates basic app with lvt new]

You: Add posts with title and content

[Agent generates posts resource]

You: Add categories and make posts belong to categories

[Agent adds categories resource and relationship]

You: Add user authentication so posts have authors

[Agent generates auth system and links posts to users]
```

The agent maintains context throughout the conversation, understanding how each piece fits together.

## Available Skills

The agent has access to the following skills:

### Core Skills (14)
1. **lvt-new-app** - Create new LiveTemplate application
2. **lvt-add-resource** - Add database-backed CRUD resource
3. **lvt-gen-auth** - Generate complete authentication system
4. **lvt-add-view** - Add standalone view (no database)
5. **lvt-customize-resource** - Modify existing resource
6. **lvt-add-relationship** - Add relationships between resources
7. **lvt-add-validation** - Add form/data validation
8. **lvt-add-middleware** - Add custom middleware
9. **lvt-customize-templates** - Modify HTML templates
10. **lvt-add-api** - Add JSON API endpoints
11. **lvt-add-websocket** - Add WebSocket support
12. **lvt-quickstart** - Fast-track common patterns
13. **lvt-troubleshoot** - Debug issues
14. **lvt-optimize** - Performance improvements

### Migration Skills (4)
1. **lvt-migration-create** - Create new migration
2. **lvt-migration-up** - Apply migrations
3. **lvt-migration-down** - Rollback migrations
4. **lvt-migration-status** - Check migration status

### Development Skills (3)
1. **lvt-run-dev** - Start development server
2. **lvt-run-tests** - Run test suite
3. **lvt-build-prod** - Build for production

## Tips for Best Results

### Be Descriptive

Instead of:
```
You: Add users
```

Try:
```
You: Add users with email, password, and profile information
```

### Ask for Relationships

```
You: Add posts that belong to users
```

The agent will understand and create the proper foreign key relationships.

### Request Features Incrementally

Start simple, then add complexity:

```
You: Create a blog
Agent: [Creates basic blog]

You: Add categories and tags
Agent: [Adds categorization]

You: Add comments
Agent: [Adds commenting system]
```

### Let the Agent Choose the Approach

The agent knows which skills to use for each task. You don't need to specify skill names.

Instead of:
```
You: Use lvt-add-resource to add posts
```

Just say:
```
You: Add posts to my blog
```

### Ask for Help

```
You: What's the best way to add user authentication?
You: How do I add real-time updates to my app?
You: I'm getting an error when I run the server
```

## Common Patterns

### Authentication Setup

```
You: I need user authentication with email and password
```

### CRUD Resource

```
You: Add products with name, price, description, and stock quantity
```

Agent will run: `lvt gen resource products name:string price:float description:text stock:int`

### Relationships

```
You: Add orders that belong to users and have many products
```

### Custom Views

```
You: Add a dashboard view that shows statistics
```

### Real-time Features

```
You: Add live updates when new messages arrive
```

## What the Agent Can Do

- **Generate Resources**: Create database-backed CRUD interfaces
- **Setup Authentication**: Complete auth system with sessions
- **Add Relationships**: One-to-many, many-to-many relationships
- **Create Migrations**: Database schema changes
- **Run Commands**: Execute lvt commands for you
- **Start Server**: Launch development server
- **Open Browser**: Navigate to your running app
- **Troubleshoot**: Help debug issues
- **Optimize**: Improve performance

## Example Session

```
You: I want to build a recipe sharing site

Agent: I'll help you build a recipe sharing site. Let me start by creating the core resources.

[Agent uses lvt-quickstart skill]
[Creates recipes resource with title, ingredients, instructions, prep_time, cook_time]
[Runs migrations]
[Starts server on port 8080]
[Opens browser to http://localhost:8080]

Agent: Your recipe site is now running! You can create, edit, and view recipes.
       Would you like me to add any additional features like:
       - User accounts and authentication
       - Recipe categories
       - Ratings and reviews
       - Photo uploads

You: Yes, add user accounts and categories

Agent: I'll add authentication and categories.

[Agent uses lvt-gen-auth skill]
[Creates users, sessions, authentication system]
[Agent uses lvt-add-resource skill]
[Creates categories resource]
[Adds relationship between recipes and categories]
[Adds relationship between recipes and users]
[Runs migrations]
[Restarts server]

Agent: ✅ Added user authentication and categories!
       - Users can now sign up and log in
       - Recipes are linked to their creators
       - Recipes can be organized into categories

       The server has been restarted. Refresh your browser to see the changes.
```

## Notes

- The agent handles migrations automatically
- Server starts on port 8080 by default
- The agent commits changes using conventional commits
- Progress is tracked in `.claude/progress/`
- All skills are documented in `.claude/skills/`

## Getting Help

If you need help with the agent:

```
You: How do I add a new feature?
You: What's the best way to structure my app?
You: I'm getting an error, can you help?
```

The agent will guide you through the process and explain what it's doing.
