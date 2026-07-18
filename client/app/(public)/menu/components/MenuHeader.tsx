import Subheading from "@/app/(landing)/homepage/SubHeading";

export default function MenuHeader() {
  return (
    <section className="py-8">
      <div className="mx-auto max-w-7xl px-6 text-center">
        <Subheading className="text-xl">Curated For Every Occasion</Subheading>

        <h1 className="mt-2 text-2xl font-semibold text-primary md:text-4xl">
          Explore Our Menu
        </h1>
      </div>
    </section>
  );
}
