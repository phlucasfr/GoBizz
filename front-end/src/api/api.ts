import { apiConfig } from "../config/apiConfig";

export interface SessionResponse {
  isValid: boolean;
  message: string;
}

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

export async function validateSession(): Promise<SessionResponse | Error> {
  try {
    const response = await fetch(`${apiConfig.baseUrl}${apiConfig.endpoints.sessions}`, {
      method: "GET",
      credentials: "include",
    });

    if (response.status === 401) {
      return { isValid: false, message: "Usuário não autorizado" };
    }

    if (!response.ok) {
      return { isValid: false, message: "Sessão inválida ou expirada" };
    }

    return { isValid: true, message: "Sessão válida" };
  } catch (error) {
    return new Error(`Erro ao validar sessão: ${error}`);
  }
}


export async function deleteSession() {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.sessions}`,
      {
        method: "DELETE",
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error("Sessão inválida");
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Erro ao remover sessão:", error);
    return null;
  }
}
