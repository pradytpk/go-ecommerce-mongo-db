package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("cannot find the product")
	ErrCantDecodeProducts = errors.New("canot find the product")
	ErrUserIDIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to ge the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
