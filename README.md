# Zynx

## Project Overview

Zynx is a dynamic messaging application designed to enhance communication and collaboration among users. With a user-friendly interface and robust features, Zynx allows individuals and groups to connect seamlessly, share information, and manage conversations effectively. 

## Tech Stack

- **Frontend**: 
  - **Next.js**: A React framework for building server-rendered applications.
  - **TypeScript**: A superset of JavaScript that adds static types, improving code quality and maintainability.

- **Backend**: 
  - **Go**: A modern programming language known for its efficiency and performance.
  - **PostgreSQL**: A powerful, open-source relational database system used for storing user data and messages.


# Project Version Control and Branch Management

Effective branch management and version control are essential for maintaining the integrity of your project and facilitating collaboration. This guide outlines the branch structure and usage in our project.

## Branch Organization

### Branch Structure

Our branching strategy follows a systematic approach to ensure clear versioning and organized development. Below is the proposed structure:

- **master**: 
  - The main stable branch that contains production-ready code.
  
- **dev**: 
  - The development branch where all features are merged before release.
  
- **release/[release-name]**: 
  - Branches created for specific releases. These branches allow for final testing and bug fixes before merging into the master branch.
  
- **feature/[release-name]/[feature-name]**: 
  - Branches dedicated to new features. Each feature branch is named according to the release it belongs to, followed by the specific feature name.
  
- **private/[feature-name]**: 
  - Branches for private development or bug fixes that are not yet ready for public view or merging into the main branches.
