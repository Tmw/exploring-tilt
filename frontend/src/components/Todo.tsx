import { ActionIcon, Card, Checkbox, Flex } from "@mantine/core";
import { type Todo } from "../api";
import { TrashIcon } from "@phosphor-icons/react";

interface TodoProps {
  item: Todo;
  onToggle: (id: string, newVal: boolean) => void;
  onDelete: (id: string) => void;
}

export function Todo({ item, onToggle, onDelete }: TodoProps) {
  const completed = item.completedAt !== null;
  const handleToggle = (event: React.ChangeEvent<HTMLInputElement>) => {
    const val = event.currentTarget.checked;
    onToggle(item.id, val);
  };

  return (
    <Card shadow="sm" padding="sm" withBorder>
      <Flex align="center" gap="md">
        <Checkbox defaultChecked={completed} onChange={handleToggle} />
        <strong
          style={{
            textDecoration: completed ? "line-through" : "none",
            flexGrow: 1,
          }}
        >
          {item.title}
        </strong>

        <ActionIcon
          variant="outline"
          size="sm"
          color="gray"
          onClick={() => onDelete(item.id)}
        >
          <TrashIcon />
        </ActionIcon>
      </Flex>
    </Card>
  );
}
