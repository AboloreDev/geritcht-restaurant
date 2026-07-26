import {
  AddToCartRequest,
  CartResponse,
  UpdateCartItemRequest,
} from "../types/cartTypes";
import { baseApi } from "./baseApi";

export const cartApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCart: builder.query<CartResponse, void>({
      query: () => ({ url: "/cart/" }),
      providesTags: ["Cart"],
    }),

    addToCart: builder.mutation<CartResponse, AddToCartRequest>({
      query: (body) => ({
        url: "/cart/",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Cart"],
    }),

    updateCartItem: builder.mutation<
      CartResponse,
      { itemId: number; body: UpdateCartItemRequest }
    >({
      query: ({ itemId, body }) => ({
        url: `/cart/${itemId}`,
        method: "PATCH",
        body,
      }),
      // optimistic update: quantity changes should feel instant,
      // not wait on a round trip before the drawer reflects it
      async onQueryStarted({ itemId, body }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          cartApi.util.updateQueryData("getCart", undefined, (draft) => {
            const item = draft.data.cart_items.find((i) => i.id === itemId);
            if (item) {
              item.quantity = body.quantity;
              item.subtotal = item.menu_item.price * body.quantity;
            }
          }),
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo(); // rollback if the request actually fails
        }
      },
      invalidatesTags: ["Cart"],
    }),

    deleteCartItem: builder.mutation<CartResponse, number>({
      query: (itemId) => ({
        url: `/cart/${itemId}`,
        method: "DELETE",
      }),
      async onQueryStarted(itemId, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          cartApi.util.updateQueryData("getCart", undefined, (draft) => {
            draft.data.cart_items = draft.data.cart_items.filter(
              (i) => i.id !== itemId,
            );
          }),
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
      invalidatesTags: ["Cart"],
    }),

    clearCart: builder.mutation<CartResponse, void>({
      query: () => ({
        url: "/cart/",
        method: "DELETE",
      }),
      invalidatesTags: ["Cart"],
    }),
  }),
});

export const {
  useGetCartQuery,
  useAddToCartMutation,
  useUpdateCartItemMutation,
  useDeleteCartItemMutation,
  useClearCartMutation,
} = cartApi;
