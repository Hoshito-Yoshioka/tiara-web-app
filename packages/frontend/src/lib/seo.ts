/**
 * サイト全体で共有する SEO 関連の定数（SEO コピーの正本モジュール）。
 * title / description / h1 サブタイトルの文言はすべてここに集約し、
 * 各 view からは PAGE_META を参照する。
 * ページごとのメタタグ生成は composables/usePageMeta.ts が担う。
 */

export const SITE_URL = 'https://tiara-hakodate.com'
export const SITE_NAME = 'Tiara'

/**
 * ホーム専用 title（主要キーワードを網羅し【公式】を明示）。
 * 要件 1.1 / 1.2: 函館・ニュークラブ・キャバクラ・ティアラ・TIARA・【公式】+ 地名補強（五稜郭）を含む。
 */
export const HOME_TITLE = '【公式】クラブ ティアラ(TIARA)｜函館・五稜郭のニュークラブ・キャバクラ'

/**
 * 下層ページ title のブランドサフィックス。
 * 要件 1.3: 「ページ固有の内容 | 函館とニュークラブを含むブランド名」の後半部。
 */
export const BRAND_TITLE = '函館のニュークラブ ティアラ【公式】'

/**
 * 旧デフォルト title。
 * usePageMeta.ts / App.vue が参照中のため後方互換として残す（HOME_TITLE への切替はタスク 1.2）。
 * @deprecated 新規コードでは HOME_TITLE / BRAND_TITLE を使用すること。
 */
export const DEFAULT_TITLE = 'Tiara（函館 ニュークラブ ティアラ）'

/**
 * サイト共通のデフォルト meta description（ホームのコピーと同一）。
 * 要件 2.2: 函館・ニュークラブ・キャバクラのキーワードと、公式サイトで実際に
 * 提供している価値（キャストプロフィール・出勤スケジュール・料金・アクセス）を含む。
 */
export const DEFAULT_DESCRIPTION =
  '函館・五稜郭のニュークラブ・キャバクラ「ティアラ(TIARA)」の公式サイト。在籍キャストのプロフィールや最新の出勤スケジュール、料金システム、店舗へのアクセスまで、公式ならではの正確な情報でご案内します。'

export const OG_IMAGE_URL = `${SITE_URL}/og-image.png`

/**
 * ページ別の固定 SEO コピー。
 * description は 100〜120 文字程度（許容 90〜125）で「函館」とページ固有の内容を含み、
 * サイトで実際に提供している価値のみを記載する（要件 2.1〜2.4）。
 */
export interface PageMetaCopy {
  /** 下層 title の固有部（ホームは undefined = HOME_TITLE を使用） */
  title?: string
  /** meta description（100〜120 文字程度） */
  description: string
  /** h1 内の可視日本語サブタイトル */
  headingSubtitle: string
}

/** 公開 6 ページ分の SEO コピーの正本。title 固有部はページ間で重複させない（要件 1.4） */
export const PAGE_META: Record<
  'home' | 'shop' | 'staff' | 'schedule' | 'price' | 'access',
  PageMetaCopy
> = {
  home: {
    description: DEFAULT_DESCRIPTION,
    headingSubtitle: '函館・五稜郭のニュークラブ ティアラ(TIARA)',
  },
  shop: {
    title: '店舗情報',
    description:
      '函館・五稜郭のニュークラブ・キャバクラ「ティアラ(TIARA)」の店舗情報ページ。住所や電話番号などご来店前に確認したい店舗の基本情報を、公式サイトならではの正確で最新の内容で分かりやすくご案内します。',
    headingSubtitle: '函館のニュークラブ ティアラの店舗情報',
  },
  staff: {
    title: 'キャスト・スタッフ紹介',
    description:
      '函館のニュークラブ・キャバクラ「ティアラ(TIARA)」の在籍キャスト・スタッフ紹介ページ。各キャストの詳しいプロフィールや写真、詳細ページへのご案内を、公式サイトならではの最新情報でご覧いただけます。',
    headingSubtitle: '函館のニュークラブ ティアラのキャスト・スタッフ紹介',
  },
  schedule: {
    title: '出勤スケジュール',
    description:
      '函館のニュークラブ・キャバクラ「ティアラ(TIARA)」の出勤スケジュールページ。店舗全体の出勤情報とキャスト個別の最新スケジュールを、公式サイトならではの正確な情報でいつでも手軽にご確認いただけます。',
    headingSubtitle: '函館のニュークラブ ティアラの出勤スケジュール',
  },
  price: {
    title: '料金システム',
    description:
      '函館のニュークラブ・キャバクラ「ティアラ(TIARA)」の料金システムページ。セット料金や各種メニューの料金体系を、すべて税込の分かりやすい表示で掲載。公式サイトならではの正確な料金情報をご案内します。',
    headingSubtitle: '函館のニュークラブ ティアラの料金システム',
  },
  access: {
    title: 'アクセス',
    description:
      '函館のニュークラブ・キャバクラ「ティアラ(TIARA)」へのアクセスページ。北海道函館市本町1-28 第5大栄ビル1F、函館市電「中央病院前」電停より徒歩約3分。地図と行き方を公式サイトがご案内します。',
    headingSubtitle: '函館のニュークラブ ティアラへのアクセス',
  },
}

/**
 * JSON-LD を <script> に埋め込む際のエスケープ。
 * 文字列中に "</script>" が含まれてもタグが閉じないよう "<" を Unicode エスケープする。
 */
export function toJsonLd(data: Record<string, unknown>): string {
  return JSON.stringify(data).replace(/</g, '\\u003c')
}
