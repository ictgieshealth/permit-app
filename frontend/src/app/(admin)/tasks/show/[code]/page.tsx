import ShowTask from "@/components/tasks/ShowTask";

interface ShowTaskPageProps {
  params: {
    code: string;
  };
}

export default function ShowTaskPage({ params }: ShowTaskPageProps) {
  return <ShowTask code={params.code} />;
}
