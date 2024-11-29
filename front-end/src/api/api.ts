export const getIpAddress = async (): Promise<string> => {
    try {
      const response = await fetch("https://api.ipify.org?format=json");
      const data = await response.json();
      return data.ip;
    } catch (error) {
      console.error("Erro ao obter o IP:", error);
      return "IP não disponível";
    }
  };