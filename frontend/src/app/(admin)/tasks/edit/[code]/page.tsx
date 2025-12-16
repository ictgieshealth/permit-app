import TaskForm from "@/components/tasks/TaskForm";

interface EditTaskPageProps {
  params: {
    code: string;
  };
}

export default function EditTaskPage({ params }: EditTaskPageProps) {
  return <TaskForm code={params.code} isEdit={true} />;
}
