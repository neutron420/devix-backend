'use client'

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import Link from 'next/link'
import { motion } from "framer-motion";


export default function FAQs() {
  const faqItems = [
    {
      id: 'item-1',
      question: 'What is Devix?',
      answer: 'Devix is a modern, real-time social platform built for software developers to share technical questions, concepts, and build logs in public.',
    },
    {
      id: 'item-2',
      question: 'What can I share on the platform?',
      answer: 'You can publish Markdown-supported articles, post real-time code updates, build logs, and ask engineering questions to the global developer community.',
    },
    {
      id: 'item-3',
      question: 'Does it support team workspaces?',
      answer: 'Yes, you can form organizations and set up collaborative developer team spaces to collaborate and share technical logs with your team.',
    },
    {
      id: 'item-4',
      question: 'Is there real-time messaging?',
      answer: 'Absolutely. Devix provides full-featured direct messaging with typing indicators, presence tracking, and instant alerts to keep your communication active.',
    },
    {
      id: 'item-5',
      question: 'How does the feed work?',
      answer: 'You can filter the global feed by trending posts and latest builds, or view a personalized stream containing updates from developers and teams you follow.',
    },
  ];


  return (
    <section className="py-16 md:py-24">
      <div className="mx-auto max-w-5xl px-6">
        <div className="grid gap-8 md:grid-cols-5 md:gap-12">
          <div className="md:col-span-2">
            <h2 className="text-foreground text-4xl font-semibold">FAQs</h2>
            <p className="text-muted-foreground mt-4 text-balance text-lg">
              Everything you need to know about Devix
            </p>
            <p className="text-muted-foreground mt-6 hidden md:block">
              Can’t find what you’re looking for? Reach out to our{' '}
              <Link
                href="#"
                className="text-primary font-medium hover:underline"
              >
                Devix support team
              </Link>{' '}
              for assistance.
            </p>
          </div>

          <div className="md:col-span-3">
            <Accordion
              type="single"
              collapsible>
              {faqItems.map((item) => (
                <AccordionItem
                  key={item.id}
                  value={item.id}
                  className="border-b border-gray-200 dark:border-gray-600">
                  <AccordionTrigger className="cursor-pointer text-base font-medium hover:no-underline">{item.question}</AccordionTrigger>
                  <AccordionContent>
                    <BlurredStagger text={item.answer} />
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          </div>

          <p className="text-muted-foreground mt-6 md:hidden">
            Can&apos;t find what you&apos;re looking for? Contact our{' '}
            <Link
              href="#"
              className="text-primary font-medium hover:underline">
              customer support team
            </Link>
          </p>
        </div>
      </div>
    </section>
  )
}

 
export const BlurredStagger = ({
  text = "built by ruixen.com",
}: {
  text: string;
}) => {
  const headingText = text;
 
  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.015,
      },
    },
  };
 
  const letterAnimation = {
    hidden: {
      opacity: 0,
      filter: "blur(10px)",
    },
    show: {
      opacity: 1,
      filter: "blur(0px)",
    },
  };
 
  return (
    <>
      <div className="w-full">
        <motion.p
          variants={container}
          initial="hidden"
          animate="show"
          className="text-base leading-relaxed break-words whitespace-normal"
        >
          {headingText.split("").map((char, index) => (
            <motion.span
              key={index}
              variants={letterAnimation}
              transition={{ duration: 0.3 }}
              className="inline-block"
            >
              {char === " " ? "\u00A0" : char}
            </motion.span>
          ))}
        </motion.p>
      </div>
    </>
  );
};