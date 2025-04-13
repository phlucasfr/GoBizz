export default function Footer() {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="bg-primary/10 py-4 text-center text-sm text-muted-foreground">
      © {currentYear} GoBizz. All rights reserved.
    </footer>
  )
}

