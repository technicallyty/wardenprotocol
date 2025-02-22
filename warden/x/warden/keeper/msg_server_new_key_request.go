package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/warden-protocol/wardenprotocol/warden/x/warden/types/v1beta2"
)

func (k msgServer) NewKeyRequest(ctx context.Context, msg *types.MsgNewKeyRequest) (*types.MsgNewKeyRequestResponse, error) {
	if err := k.assertIntentAuthority(msg.Authority); err != nil {
		return nil, err
	}

	creator := k.intentKeeper.GetActionCreator(ctx)

	if _, err := k.SpacesKeeper.Get(ctx, msg.SpaceId); err != nil {
		return nil, err
	}

	keychain, err := k.keychains.Get(ctx, msg.KeychainId)
	if err != nil {
		return nil, err
	}

	if keychain.Fees != nil {
		err := k.bankKeeper.SendCoins(
			ctx,
			sdk.MustAccAddressFromBech32(creator),
			keychain.AccAddress(),
			sdk.NewCoins(sdk.NewInt64Coin("uward", keychain.Fees.KeyReq)),
		)
		if err != nil {
			return nil, err
		}
	}

	req := &types.KeyRequest{
		Creator:    creator,
		SpaceId:    msg.SpaceId,
		KeychainId: msg.KeychainId,
		KeyType:    msg.KeyType,
		Status:     types.KeyRequestStatus_KEY_REQUEST_STATUS_PENDING,
		IntentId:   msg.IntentId,
	}

	id, err := k.keyRequests.Append(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgNewKeyRequestResponse{
		Id: id,
	}, nil
}
