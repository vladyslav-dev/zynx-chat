export const validateEnvVariables = () => {
    const envVars = {
      SERVER_URL: process.env.REACT_APP_SERVER_URL,
      API_URL: process.env.REACT_APP_API_URL,
      WEBSOCKET_URL: process.env.REACT_APP_WEBSOCKET_URL,
    };
  
    const missingEnvVars = Object.entries(envVars)
      .filter(([_, value]) => !value)
      .map(([key]) => key);
  
    return missingEnvVars.length === 0;
};