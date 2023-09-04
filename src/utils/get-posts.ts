import Parser from 'rss-parser';
import { z } from 'zod';

export const getPosts = async () => {
  const url = 'https://lifeofdev.com/rss.xml';
  const rss = await new Parser().parseURL(url);

  // Sort the items by pubDate before validation
  const sortedItems = rss.items.sort((a, b) => {
    if (typeof a.pubDate === 'string' && typeof b.pubDate === 'string') {
      return new Date(b.pubDate).valueOf() - new Date(a.pubDate).valueOf();
    }
    return 0;
  });

  // Take only the first 5 items after sorting
  const first5Items = sortedItems.slice(0, 5);

  const posts = z
    .array(
      z.object({
        title: z.string(),
        link: z.string(),
        pubDate: z.string().transform((date) => new Date(date).valueOf()),
      }),
    )
    .safeParse(first5Items);

  if (!posts.success) throw new Error(posts.error.message);

  return posts.data;
};
