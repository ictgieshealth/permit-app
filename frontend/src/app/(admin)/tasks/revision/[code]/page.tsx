import RevisionForm from "@/components/tasks/RevisionForm";

interface RevisionTaskPageProps {
  params: {
    code: string;
  };
}

export default function RevisionTaskPage({ params }: RevisionTaskPageProps) {
  return <RevisionForm code={params.code} />;
}
