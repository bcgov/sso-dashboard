"""allow null in offline sessions

Revision ID: 7c002818cb32
Revises: 5bfc2a71a71f
Create Date: 2024-07-29 16:29:07.464223

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '7c002818cb32'
down_revision = '5bfc2a71a71f'
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.alter_column('client_sessions','offline_sessions', existing_type=sa.INTEGER(), nullable=True)


def downgrade() -> None:
    op.alter_column('client_sessions','offline_sessions', existing_type=sa.INTEGER(), nullable=False)
