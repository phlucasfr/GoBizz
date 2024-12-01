import { useParams } from "@solidjs/router";

const Home = () => {
  const { id } = useParams();

  return (
    <div>
      <h1>Home - ID: {id}</h1>
    </div>
  );
};

export default Home;
