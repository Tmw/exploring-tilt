import { Container } from "@mantine/core";

export function Page(props: React.PropsWithChildren) {
  return (
    <Container strategy="grid" size="md">
      {props.children}
    </Container>
  );
}
