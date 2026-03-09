import { Injectable, NotFoundException } from '@nestjs/common';
import { v4 as uuidv4 } from 'uuid';
import { ProductsService } from '../products/products.service';
import { OrdersRepository } from './orders.repository';
import { CreateOrderDto } from './dto/create-order.dto';
import { Order, OrderItem } from './order.entity';

@Injectable()
export class OrdersService {
  constructor(
    private readonly ordersRepository: OrdersRepository,
    private readonly productsService: ProductsService,
  ) {}

  create(dto: CreateOrderDto): Order {
    const items: OrderItem[] = dto.items.map(({ productId, quantity }) => {
      const product = this.productsService.findById(productId);
      const unitVatInCents = Math.round(product.priceInCents * product.vatRate);

      return {
        productId: product.id,
        productName: product.name,
        quantity,
        unitPriceInCents: product.priceInCents,
        unitVatInCents,
        vatRate: product.vatRate * 100,
      };
    });

    const totalPriceInCents = items.reduce(
      (sum, item) => sum + item.unitPriceInCents * item.quantity,
      0,
    );

    const totalVatInCents = items.reduce(
      (sum, item) => sum + item.unitVatInCents * item.quantity,
      0,
    );

    const order: Order = {
      id: uuidv4(),
      items,
      totalPriceInCents,
      totalVatInCents,
      createdAt: new Date(),
    };

    return this.ordersRepository.save(order);
  }

  findById(id: string): Order {
    const order = this.ordersRepository.findById(id);
    if (!order) {
      throw new NotFoundException(`Order with id "${id}" not found`);
    }
    return order;
  }

  findAll(): Order[] {
    return this.ordersRepository.findAll();
  }
}
