import { motion } from "framer-motion";

interface ItemListProps {
  title: string;
  price: string;
  tags: string;
}

export default function ItemList({ title, price, tags }: ItemListProps) {
  return (
    <motion.div
      whileHover={{ x: 6 }}
      transition={{ duration: 0.25 }}
      className="group space-y-2"
    >
      <div className="flex items-center gap-4">
        <h4 className="font-heading text-xl text-primary transition-colors duration-300 group-hover:text-primary-deep">
          {title}
        </h4>

        <div className="h-px flex-1 bg-gradient-to-r from-primary/70 to-transparent" />

        <p className="font-heading text-lg font-semibold text-white">{price}</p>
      </div>

      <p className="text-sm tracking-wide text-text-muted">{tags}</p>
    </motion.div>
  );
}
