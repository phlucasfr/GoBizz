import SignUp from "./SignUp";

const Welcome = () => {
  return (
    <div class="flex h-[84vh] overflow-hidden rounded-lg">
      <div class="md:flex md:w-2/3 lg:w-3/4 rounded-lg text-white items-center justify-center">
        <div class="hidden  md:block text-center max-w-xl">
          <h1 class="text-4xl font-bold mb-6 text-black">
            Bem-vindo à Nossa Plataforma
          </h1>
          <p class="text-xl text-black">
            Transforme sua gestão empresarial com ferramentas inteligentes e
            recursos completos.
          </p>
        </div>

        <div class="absolute top-0 right-0 bottom-0 w-full md:w-1/3 lg:w-1/4 flex items-center justify-center p-6">
          <SignUp />
        </div>
      </div>
    </div>
  );
};

export default Welcome;
