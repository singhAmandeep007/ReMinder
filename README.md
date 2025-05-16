# ReMinder - A Production Grade Reminder App 🚀

A fully-featured reminder application showcasing modern web development practices, backend integration, and testing methodologies. This project is designed to be a reference for building scalable and maintainable applications using React, TypeScript, NestJS, Golang, CI/CD implementation and various testing frameworks.

## 🎯 Key Features

- **Robust Frontend Architecture**: Built with React, TypeScript, and Redux Toolkit
- **Modern UI**: Responsive design using RadixUI + Tailwind CSS with dark/light themes
- **Internationalization**: Multi-language support with i18n
- **Comprehensive Testing**:
  - E2E Testing with Cypress and Playwright
  - Unit Testing with Jest and React Testing Library
  - API Mocking with MSW and Mirage.js
- **CI/CD Pipeline**: Automated testing, building, and deployment
- **Best Practices**: SOLID principles, TDD approach, Clean Architecture

## 🛠️ Tech Stack

### Frontend
- React
- TypeScript
- Redux Toolkit + RTK Query
- Shadcn UI
- Tailwind CSS
- React Router
- I18n for localization

### Backend
- [NestJS](/backend/nestjs-server/README.md)
- [Golang](/backend/gin-server/README.md)

### Testing
- Jest & React Testing Library
- Cypress
- Playwright
- MSW (Mock Service Worker) + MSW Data
- Mirage.js
- Coverage reports with Istanbul + Mocha

### Development Tools
- Storybook
- ESLint & Prettier
- Husky & lint-staged
- Commitizen
- TypeScript
- Webpack

## 📚 Quick Start

```bash
# Clone the repository
git clone https://github.com/singhAmandeep007/ReMinder.git

# Install dependencies - Requires Node.js v20.11.1
npm install

# Start development server
npm start

# Run tests
npm test:unit
```

## 📖 Available Scripts

- `npm start` - Start development server
- `npm run analyze` - Analyze bundle size
- `npm run preview` - Run the production build locally - API calls will be proxied to Mock Server
- `npm run lint` - Run formatting and linting
- `npm run deploy` - Deploy the app to github pages along with docs
- `npm run storybook:start` - Run Storybook locally
- `npm run commit:init` - Initialize a commit process for staged files

- [Testing scripts](./testing.md)

## 🌟 Best Practices

- Test-Driven Development (TDD)
- SOLID Principles
- Component-Driven Development
- Atomic Design Principles
- Clean Code Architecture
- Conventional Commits

## 🤝 Contributing

Contributions are welcome!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
