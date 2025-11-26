# LiveTemplate Agent Usage Guide

This guide explains how to use the LiveTemplate assistant agent in Claude Code to create and develop applications.

## Getting Started

### 1. Create a New Project

```bash
# Create a new LiveTemplate app
lvt new myapp

# Navigate to the project
cd myapp
```

### 2. Open in Claude Code

```bash
# Open the project in Claude Code CLI
claude
```

The agent will be automatically available in your Claude Code session.

Note: If you don't have the Claude CLI installed, you can open the directory in your preferred IDE with the Claude Code extension.

### 3. Start a Conversation

Simply ask the agent what you want to build. The agent understands natural language and will guide you through the process.

## Example Workflows

### Quick Start (Simplest Approach)

```
You: I want to build a blog with posts
```

The agent will:
- Use the `lvt-quickstart` skill
- Generate a posts resource with appropriate fields
- Run migrations
- Start the server
- Open your browser to the running app

### Full Stack Example

```
You: I want to build a task management system with:
- Users with authentication
- Projects that belong to users
- Tasks that belong to projects
- Due dates and priorities
```

The agent will:
- Use `lvt-gen-auth` to set up authentication
- Use `lvt-add-resource` for projects with user relationships
- Use `lvt-add-resource` for tasks with project relationships
- Configure the appropriate field types and relationships
- Run all migrations
- Start the server

### Step-by-Step Example

```
You: Let's add a blog to my app

Agent: I'll help you add a blog. Let me create a posts resource.

[Agent creates posts with title, content, published_at fields]

You: Can you add categories and tags?

Agent: I'll add those resources with relationships.

[Agent adds categories, tags, and junction tables]

You: Add author attribution to posts

Agent: I'll add an author relationship to the posts resource.

[Agent modifies posts to include author_id field]
```

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
