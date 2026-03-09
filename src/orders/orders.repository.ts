import { Injectable } from '@nestjs/common';
import { Order } from './order.entity';

@Injectable()
export class OrdersRepository {
  private readonly orders = new Map<string, Order>();

  save(order: Order): Order {
    this.orders.set(order.id, order);
    return order;
  }

  findById(id: string): Order | undefined {
    return this.orders.get(id);
  }

  findAll(): Order[] {
    return Array.from(this.orders.values());
  }
}
