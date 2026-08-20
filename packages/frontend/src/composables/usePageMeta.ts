import { computed, toValue, type MaybeRefOrGetter } from 'vue'
import { useHead, useSeoMeta } from '@unhead/vue'
import { useRoute } from 'vue-router'
import { SITE_URL, SITE_NAME, DEFAULT_TITLE, DEFAULT_DESCRIPTION, OG_IMAGE_URL } from '@/lib/seo'

interface PageMetaInput {
  /** ページ固有のタイトル。省略時はサイト共通タイトルのみ */
  title?: MaybeRefOrGetter<string | undefined>
  /** ページ固有の description。省略時はサイト共通 description */
  description?: MaybeRefOrGetter<string | undefined>
}

/**
 * 公開ページの title / description / canonical / OGP / Twitter Card を設定する。
 * SSG ビルド時は各ページの HTML に静的に埋め込まれる。
 * title・description にはゲッターを渡せるため、API 取得後の動的な値にも追従する。
 */
export function usePageMeta(input: PageMetaInput = {}) {
  const route = useRoute()

  const canonicalUrl = computed(() => `${SITE_URL}${route.path}`)
  const title = computed(() => {
    const pageTitle = toValue(input.title)
    return pageTitle ? `${pageTitle} | ${DEFAULT_TITLE}` : DEFAULT_TITLE
  })
  const description = computed(() => toValue(input.description) || DEFAULT_DESCRIPTION)

  useHead({
    link: [{ rel: 'canonical', href: canonicalUrl }],
  })

  useSeoMeta({
    title,
    description,
    ogTitle: title,
    ogDescription: description,
    ogUrl: canonicalUrl,
    ogSiteName: SITE_NAME,
    ogType: 'website',
    ogLocale: 'ja_JP',
    ogImage: OG_IMAGE_URL,
    twitterCard: 'summary_large_image',
  })
}
