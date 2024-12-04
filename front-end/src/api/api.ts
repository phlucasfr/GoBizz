import { apiConfig } from "../config/apiConfig";

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

export async function loginCompany(email: string, password: string) {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.companyLogin}`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      }
    );

    if (response.status === 401) {
      return { success: false, message: "Invalid password or email" };
    }

    if (response.status === 404) {
      return { success: false, message: "Empresa não encontrada" };
    }

    if (!response.ok) {
      const errorData = await response.json();
      return {
        success: false,
        message: errorData.error || "Erro desconhecido ao fazer login",
      };
    }

    return { success: true, message: "Login efetuado com sucesso" };
  } catch (error) {
    console.error("Erro ao fazer login:", error);
    return {
      success: false,
      message: `Erro ao fazer login: ${
        error instanceof Error ? error.message : error
      }`,
    };
  }
}

export async function validateSession() {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.sessions}`,
      {
        method: "GET",
        credentials: "include",
      }
    );

    if (response.status === 401) {
      return { isValid: false, message: "Usuário não autorizado" };
    }

    if (response.status === 403) {
      return { isValid: false, message: "Empresa desativada" };
    }

    if (!response.ok) {
      const errorData = await response.json();
      return {
        isValid: false,
        message: errorData?.error || "Sessão inválida ou expirada",
      };
    }

    return { isValid: true, message: "Sessão válida" };
  } catch (error) {
    console.error("Erro ao validar sessão:", error);
    return {
      isValid: false,
      message: `Erro ao validar sessão: ${error}`,
    };
  }
}

export async function sendVerificationEmail(email: string) {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.emailVerification}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      return {
        success: false,
        message: errorData?.error || "Erro ao enviar o e-mail de verificação.",
      };
    }

    return {
      success: true,
      message: "E-mail de verificação enviado com sucesso.",
    };
  } catch (error) {
    console.error("Erro ao enviar e-mail de verificação:", error);
    return {
      success: false,
      message: "Erro inesperado ao tentar enviar o e-mail de verificação.",
    };
  }
}

export async function sendPasswordRecoveryEmail(email: string) {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.recovery}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      }
    );

    if (response.status === 404) {
      return { success: false, message: "Empresa não encontrada" };
    }

    if (!response.ok) {
      const errorData = await response.json();
      return {
        success: false,
        message: errorData?.error || "Erro ao enviar o e-mail de recuperação.",
      };
    }

    return {
      success: true,
      message: "E-mail de recuperação enviado com sucesso.",
    };
  } catch (error) {
    console.error("Erro ao enviar e-mail de recuperação:", error);
    return {
      success: false,
      message: "Erro inesperado ao tentar enviar o e-mail de recuperação.",
    };
  }
}

export async function resetPassword(token: string, password: string) {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.resetPassword}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ token, password }),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      return {
        success: false,
        message: errorData?.error || "Erro ao alterar a senha.",
      };
    }

    return {
      success: true,
      message: "Senha alterada com sucesso.",
    };
  } catch (error) {
    console.error("Erro ao alterar a senha:", error);
    return {
      success: false,
      message: "Erro inesperado ao tentar enviar o e-mail de recuperação.",
    };
  }
}

export async function VerifyCompanyByEmail(token: string) {
  try {
    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.emailVerification}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ token }),
      }
    );

    const data = await response.json();
    if (!response.ok) {
      return {
        success: false,
        message: data?.error || "Erro ao verificar empresa.",
      };
    }

    return {
      success: true,
      message: data?.message || "Empresa verificada com sucesso.",
    };
  } catch (error) {
    console.error("Erro ao verificar a empresa:", error);
    return {
      success: false,
      message: "Erro inesperado ao tentar verificar a empresa.",
    };
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
      const errorData = await response.json();
      return {
        success: false,
        message: errorData?.error || "Sessão inválida ou expirada",
      };
    }

    return { success: true, message: "Sessão removida com sucesso" };
  } catch (error) {
    console.error("Erro ao remover sessão:", error);
    return {
      success: false,
      message: `Erro ao remover sessão: ${error}`,
    };
  }
}
