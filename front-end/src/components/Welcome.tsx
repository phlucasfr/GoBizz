import SignUp from "./SignUp";

const Welcome = () => {
  return (
    <div class="flex flex-col md:flex-row h-screen">
      <div class="flex flex-col items-center justify-center w-full md:w-2/3 p-6 bg-gray-100">
        <div class="text-center max-w-xl">
          <h1 class="text-2xl md:text-4xl font-bold mb-6 text-black">
            Bem-vindo à Nossa Plataforma
          </h1>
          <p class="text-base md:text-xl text-black">
            Transforme sua gestão empresarial com ferramentas inteligentes e
            recursos completos.
          </p>
        </div>
      </div>

      <div class="w-full md:w-1/3 flex items-center justify-center p-6 overflow-auto">
        <SignUp />
      </div>
    </div>
  );
};

export default Welcome;
